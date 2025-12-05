"use client"

import { useState, useEffect } from "react"
import { Plus, Edit, Trash2, Upload } from "lucide-react"
import { Button } from "@/components/ui/button"
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table"
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog"
import { Input } from "@/components/ui/input"
import { Label } from "@/components/ui/label"
import { Textarea } from "@/components/ui/textarea"
import { Switch } from "@/components/ui/switch"
import { useToast } from "@/hooks/use-toast"
import api from "@/lib/api"
import type { Channel } from "@/types"

export default function ChannelsPage() {
  const [channels, setChannels] = useState<Channel[]>([])
  const [loading, setLoading] = useState(true)
  const [isCreateOpen, setIsCreateOpen] = useState(false)
  const [isBatchOpen, setIsBatchOpen] = useState(false)
  const [editingChannel, setEditingChannel] = useState<Channel | null>(null)
  const { toast } = useToast()

  const [formData, setFormData] = useState({
    id: 0,
    channel_id: "",
    display_name: "",
    is_active: 1,
    regexp: "",
    category: "",
    area: "CN",
    logo_url: "",
    timezone: "Asia/Shanghai",
  })

  const [batchData, setBatchData] = useState("")

  useEffect(() => {
    loadChannels()
  }, [])

  const loadChannels = async () => {
    try {
      setLoading(true)
      const response = await api.get("/admin/channels")
      setChannels(response.data || [])
    } catch (error: any) {
      toast({
        variant: "destructive",
        title: "加载失败",
        description: error.response?.data?.error || "无法加载频道列表",
      })
    } finally {
      setLoading(false)
    }
  }

  const handleCreate = async () => {
    try {
      await api.post("/admin/channels", formData)
      toast({
        title: "创建成功",
        description: "频道已成功创建",
      })
      setIsCreateOpen(false)
      resetForm()
      loadChannels()
    } catch (error: any) {
      toast({
        variant: "destructive",
        title: "创建失败",
        description: error.response?.data?.error || "创建频道失败",
      })
    }
  }

  const handleUpdate = async () => {
    if (!editingChannel) return

    try {
      await api.put(`/admin/channels/${editingChannel.channel_id}`, formData)
      toast({
        title: "更新成功",
        description: "频道信息已更新",
      })
      setEditingChannel(null)
      loadChannels()
      resetForm()
    } catch (error: any) {
      toast({
        variant: "destructive",
        title: "更新失败",
        description: error.response?.data?.error || "更新频道失败",
      })
    }
  }

  const handleDelete = async (channelId: string) => {
    if (!confirm("确定要删除这个频道吗?")) return

    try {
      await api.delete(`/admin/channels/${channelId}`)
      toast({
        title: "删除成功",
        description: "频道已被删除",
      })
      loadChannels()
    } catch (error: any) {
      toast({
        variant: "destructive",
        title: "删除失败",
        description: error.response?.data?.error || "删除频道失败",
      })
    }
  }

  const handleBatchCreate = async () => {
    try {
      const lines = batchData
        .trim()
        .split("\n")
        .filter((line) => line.trim())
      const channels = lines.map((line) => {
        const [
          channel_id,
          display_name,
          category = "",
          area = "CN",
          logo_url = "",
          timezone = "Asia/Shanghai",
        ] = line.split(",").map((s) => s.trim())
        return { channel_id, display_name, category, area, logo_url, timezone }
      })

      await api.post("/admin/channels/batch", { channels })
      toast({
        title: "批量创建成功",
        description: `成功创建 ${channels.length} 个频道`,
      })
      setIsBatchOpen(false)
      setBatchData("")
      loadChannels()
    } catch (error: any) {
      toast({
        variant: "destructive",
        title: "批量创建失败",
        description: error.response?.data?.error || "批量创建频道失败",
      })
    }
  }

  const openEdit = (channel: Channel) => {
    setEditingChannel(channel)
    setFormData({
      id: channel.id,
      is_active: channel.is_active,
      channel_id: channel.channel_id,
      display_name: channel.display_name,
      regexp: channel.regexp || "",
      category: channel.category || "",
      area: channel.area || "CN",
      logo_url: channel.logo_url || "",
      timezone: channel.timezone || "Asia/Shanghai",
    })
  }

  const resetForm = () => {
    setFormData({
      id: 0,
      is_active: 1,
      channel_id: "",
      display_name: "",
      category: "",
      regexp: "",
      area: "CN",
      logo_url: "",
      timezone: "Asia/Shanghai",
    })
  }

  return (
    <div className="space-y-6">
      <div className="flex justify-between items-center">
        <div>
          <h1 className="text-3xl font-bold">频道管理</h1>
          <p className="text-gray-500 mt-1">管理所有频道信息</p>
        </div>
        <div className="flex gap-2">
          <Button onClick={() => setIsBatchOpen(true)} variant="outline">
            <Upload className="h-4 w-4 mr-2" />
            批量导入
          </Button>
          <Button onClick={() => setIsCreateOpen(true)}>
            <Plus className="h-4 w-4 mr-2" />
            新建频道
          </Button>
        </div>
      </div>

      {loading ? (
        <div className="text-center py-12">加载中...</div>
      ) : (
        <div className="bg-white rounded-lg shadow">
          <Table>
            <TableHeader>
              <TableRow>
                <TableHead>频道ID</TableHead>
                <TableHead>名称</TableHead>
                <TableHead>匹配正则</TableHead>
                <TableHead>分类</TableHead>
                <TableHead>地区</TableHead>
                <TableHead>时区</TableHead>
                <TableHead>状态</TableHead>
                <TableHead className="text-right">操作</TableHead>
              </TableRow>
            </TableHeader>
            <TableBody>
              {channels.map((channel) => (
                <TableRow key={channel.id}>
                  <TableCell className="font-mono text-sm">
                    {channel.channel_id}
                  </TableCell>
                  <TableCell className="font-medium">
                    {channel.display_name}
                  </TableCell>
                  <TableCell className="font-medium">
                    {channel.regexp || "-"}
                  </TableCell>
                  <TableCell>{channel.category || "-"}</TableCell>
                  <TableCell>{channel.area || "-"}</TableCell>
                  <TableCell className="text-sm">
                    {channel.timezone || "-"}
                  </TableCell>
                  <TableCell>
                    <span
                      className={`px-2 py-1 rounded text-xs ${
                        channel.is_active
                          ? "bg-green-100 text-green-800"
                          : "bg-gray-100 text-gray-800"
                      }`}
                    >
                      {channel.is_active ? "活跃" : "禁用"}
                    </span>
                  </TableCell>
                  <TableCell className="text-right">
                    <Button
                      variant="ghost"
                      size="sm"
                      onClick={() => openEdit(channel)}
                    >
                      <Edit className="h-4 w-4" />
                    </Button>
                    <Button
                      variant="ghost"
                      size="sm"
                      onClick={() => handleDelete(channel.channel_id)}
                    >
                      <Trash2 className="h-4 w-4 text-red-500" />
                    </Button>
                  </TableCell>
                </TableRow>
              ))}
            </TableBody>
          </Table>
        </div>
      )}

      <Dialog
        open={isCreateOpen || !!editingChannel}
        onOpenChange={(open) => {
          if (!open) {
            setIsCreateOpen(false)
            setEditingChannel(null)
            resetForm()
          }
        }}
      >
        <DialogContent className="max-w-md">
          <DialogHeader>
            <DialogTitle>
              {editingChannel ? "编辑频道" : "新建频道"}
            </DialogTitle>
            <DialogDescription>
              {editingChannel ? "修改频道信息" : "添加新的频道"}
            </DialogDescription>
          </DialogHeader>
          <div className="space-y-4">
            <div className="space-y-2">
              <Label htmlFor="channel_id">频道ID *</Label>
              <Input
                id="channel_id"
                value={formData.channel_id}
                onChange={(e) =>
                  setFormData({ ...formData, channel_id: e.target.value })
                }
                placeholder="例如: CCTV1"
                disabled={!!editingChannel}
              />
            </div>
            <div className="space-y-2">
              <Label htmlFor="display_name">显示名称 *</Label>
              <Input
                id="display_name"
                value={formData.display_name}
                onChange={(e) =>
                  setFormData({ ...formData, display_name: e.target.value })
                }
                placeholder="例如: CCTV-1综合"
              />
            </div>
            <div className="space-y-2">
              <Label htmlFor="regexp">匹配正则</Label>
              <Input
                id="regexp"
                value={formData.regexp}
                onChange={(e) =>
                  setFormData({ ...formData, regexp: e.target.value })
                }
                placeholder="例如: ^CCTV-?1"
              />
            </div>
            <div className="space-y-2">
              <Label htmlFor="is_active">是否启用</Label>
              <Switch
                checked={formData.is_active === 1}
                onCheckedChange={(checked) =>
                  setFormData({ ...formData, is_active: checked ? 1 : 0 })
                }
              ></Switch>
            </div>
            <div className="space-y-2">
              <Label htmlFor="category">分类</Label>
              <Input
                id="category"
                value={formData.category}
                onChange={(e) =>
                  setFormData({ ...formData, category: e.target.value })
                }
                placeholder="例如: 央视"
              />
            </div>
            <div className="space-y-2">
              <Label htmlFor="area">地区</Label>
              <Input
                id="area"
                value={formData.area}
                onChange={(e) =>
                  setFormData({ ...formData, area: e.target.value })
                }
                placeholder="例如: CN"
              />
            </div>
            <div className="space-y-2">
              <Label htmlFor="logo_url">Logo URL</Label>
              <Input
                id="logo_url"
                value={formData.logo_url}
                onChange={(e) =>
                  setFormData({ ...formData, logo_url: e.target.value })
                }
                placeholder="https://..."
              />
            </div>
            <div className="space-y-2">
              <Label htmlFor="timezone">时区</Label>
              <Input
                id="timezone"
                value={formData.timezone}
                onChange={(e) =>
                  setFormData({ ...formData, timezone: e.target.value })
                }
                placeholder="Asia/Shanghai"
              />
            </div>
          </div>
          <DialogFooter>
            <Button
              variant="outline"
              onClick={() => {
                setIsCreateOpen(false)
                setEditingChannel(null)
                resetForm()
              }}
            >
              取消
            </Button>
            <Button onClick={editingChannel ? handleUpdate : handleCreate}>
              {editingChannel ? "更新" : "创建"}
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>

      <Dialog open={isBatchOpen} onOpenChange={setIsBatchOpen}>
        <DialogContent className="max-w-2xl">
          <DialogHeader>
            <DialogTitle>批量导入频道</DialogTitle>
            <DialogDescription>
              每行一个频道,格式: 频道ID,显示名称,分类,地区,Logo URL,时区
            </DialogDescription>
          </DialogHeader>
          <div>
            <Textarea
              value={batchData}
              onChange={(e) => setBatchData(e.target.value)}
              placeholder="cctv1,CCTV-1综合,央视,CN,,Asia/Shanghai&#10;cctv2,CCTV-2财经,央视,CN,,Asia/Shanghai"
              rows={10}
              className="font-mono text-sm"
            />
          </div>
          <DialogFooter>
            <Button variant="outline" onClick={() => setIsBatchOpen(false)}>
              取消
            </Button>
            <Button onClick={handleBatchCreate}>导入</Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>
    </div>
  )
}
