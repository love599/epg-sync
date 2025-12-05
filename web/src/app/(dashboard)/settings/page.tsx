"use client"

import { useState } from "react"
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import { Label } from "@/components/ui/label"
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card"
import { useToast } from "@/hooks/use-toast"
import { authApi } from "@/lib/api"
import { useAuthStore } from "@/store/authStore"
import { Lock, User } from "lucide-react"

export default function SettingsPage() {
  const { toast } = useToast()
  const user = useAuthStore((state) => state.user)
  const [loading, setLoading] = useState(false)
  const [formData, setFormData] = useState({
    oldPassword: "",
    newPassword: "",
    confirmPassword: "",
  })

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()

    if (!formData.oldPassword || !formData.newPassword || !formData.confirmPassword) {
      toast({
        variant: "destructive",
        title: "请填写完整信息",
        description: "所有密码字段都不能为空",
      })
      return
    }

    if (formData.newPassword.length < 8) {
      toast({
        variant: "destructive",
        title: "密码太短",
        description: "新密码长度至少为8个字符",
      })
      return
    }

    if (formData.newPassword !== formData.confirmPassword) {
      toast({
        variant: "destructive",
        title: "密码不一致",
        description: "新密码和确认密码不匹配",
      })
      return
    }

    try {
      setLoading(true)
      await authApi.changePassword(formData.oldPassword, formData.newPassword)
      toast({
        title: "修改成功",
        description: "密码已成功修改",
      })
      setFormData({
        oldPassword: "",
        newPassword: "",
        confirmPassword: "",
      })
    } catch (error: unknown) {
      toast({
        variant: "destructive",
        title: "修改失败",
        description:
          error instanceof Error ? error.message : "旧密码错误或服务器异常",
      })
    } finally {
      setLoading(false)
    }
  }

  return (
    <div className="space-y-6">
      <div>
        <h1 className="text-3xl font-bold text-gray-900">账户设置</h1>
        <p className="mt-2 text-sm text-gray-600">
          管理您的账户信息和安全设置
        </p>
      </div>

      <div className="grid gap-6 lg:grid-cols-2">
        <Card>
          <CardHeader>
            <CardTitle className="flex items-center">
              <User className="mr-2 h-5 w-5" />
              账户信息
            </CardTitle>
            <CardDescription>您的基本账户信息</CardDescription>
          </CardHeader>
          <CardContent className="space-y-4">
            <div>
              <Label className="text-sm text-gray-500">用户名</Label>
              <p className="text-base font-medium text-gray-900">
                {user?.username}
              </p>
            </div>
            <div>
              <Label className="text-sm text-gray-500">邮箱</Label>
              <p className="text-base font-medium text-gray-900">
                {user?.email || "未设置"}
              </p>
            </div>
            <div>
              <Label className="text-sm text-gray-500">角色</Label>
              <p className="text-base font-medium text-gray-900">
                {user?.role}
              </p>
            </div>
          </CardContent>
        </Card>

        <Card>
          <CardHeader>
            <CardTitle className="flex items-center">
              <Lock className="mr-2 h-5 w-5" />
              修改密码
            </CardTitle>
            <CardDescription>更新您的登录密码</CardDescription>
          </CardHeader>
          <CardContent>
            <form onSubmit={handleSubmit} className="space-y-4">
              <div className="space-y-2">
                <Label htmlFor="oldPassword">旧密码</Label>
                <Input
                  id="oldPassword"
                  type="password"
                  placeholder="请输入旧密码"
                  value={formData.oldPassword}
                  onChange={(e) =>
                    setFormData({ ...formData, oldPassword: e.target.value })
                  }
                  disabled={loading}
                  autoComplete="current-password"
                />
              </div>
              <div className="space-y-2">
                <Label htmlFor="newPassword">新密码</Label>
                <Input
                  id="newPassword"
                  type="password"
                  placeholder="请输入新密码 (至少6个字符)"
                  value={formData.newPassword}
                  onChange={(e) =>
                    setFormData({ ...formData, newPassword: e.target.value })
                  }
                  disabled={loading}
                  autoComplete="new-password"
                />
              </div>
              <div className="space-y-2">
                <Label htmlFor="confirmPassword">确认新密码</Label>
                <Input
                  id="confirmPassword"
                  type="password"
                  placeholder="请再次输入新密码"
                  value={formData.confirmPassword}
                  onChange={(e) =>
                    setFormData({
                      ...formData,
                      confirmPassword: e.target.value,
                    })
                  }
                  disabled={loading}
                  autoComplete="new-password"
                />
              </div>
              <Button type="submit" disabled={loading} className="w-full">
                {loading ? "修改中..." : "修改密码"}
              </Button>
            </form>
          </CardContent>
        </Card>
      </div>
    </div>
  )
}
