export interface Channel {
  id: number;
  channel_id: string;
  display_name: string;
  category?: string;
  regexp?: string;
  area?: string;
  logo_url?: string;
  timezone?: string;
  is_active: number;
  created_at: string;
  updated_at: string;
}

export interface ChannelMapping {
  id: number;
  canonical_id: string;
  provider_id: string;
  provider_channel_id: string;
  provider_channel_name?: string;
  confidence: number;
  is_verified: number;
  created_at: string;
  updated_at: string;
}

export interface Program {
  id: number;
  channel_id: string;
  title: string;
  start_time: string;
  end_time: string;
  description?: string;
  provider_id: string;
  original_timezone: string;
  created_at: string;
  updated_at: string;
}

export interface User {
  id: number;
  username: string;
  email?: string;
  role: string;
}

export interface LoginRequest {
  username: string;
  password: string;
}

export interface LoginResponse {
  message: string;
  data: {
    token: string;
    user: User;
  };
}

export interface ApiResponse<T = any> {
  data?: T;
  message?: string;
  error?: string;
  code?: string;
}
