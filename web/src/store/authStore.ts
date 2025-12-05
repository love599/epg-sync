import { create } from 'zustand';
import { persist, createJSONStorage } from 'zustand/middleware';
import api from '@/lib/api';

interface User {
  id: number;
  username: string;
  email: string;
  role: string;
}

interface AuthState {
  user: User | null;
  token: string | null;
  isHydrated: boolean;
  login: (username: string, password: string) => Promise<void>;
  logout: () => void;
  setHydrated: () => void;
}

export const useAuthStore = create<AuthState>()(
  persist(
    (set) => ({
      user: null,
      token: null,
      isHydrated: false,

      login: async (username: string, password: string) => {
        const response = await api.post('/auth/login', {
          username,
          password,
        });

        const { token, user } = response.data;

        
        api.defaults.headers.common['Authorization'] = `Bearer ${token}`;

        set({ token, user });
      },

      logout: () => {
        delete api.defaults.headers.common['Authorization'];
        set({ token: null, user: null });
      },

      setHydrated: () => {
        set({ isHydrated: true });
      },
    }),
    {
      name: 'auth-storage',
      storage: createJSONStorage(() => localStorage),
      onRehydrateStorage: () => (state) => {
        if (state?.token) {
          api.defaults.headers.common['Authorization'] = `Bearer ${state.token}`;
        }
        state?.setHydrated();
      },
    }
  )
);
