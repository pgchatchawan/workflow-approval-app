"use client";

export default function ApprovalModal({
  open,
  title,
  confirmText,
  confirmVariant,
  reason,
  setReason,
  onClose,
  onConfirm,
  loading,
}: {
  open: boolean;
  title: string;
  confirmText: string;
  confirmVariant: "approve" | "reject";
  reason: string;
  setReason: (v: string) => void;
  onClose: () => void;
  onConfirm: () => void;
  loading: boolean;
}) {
  if (!open) return null;

  const confirmCls =
    confirmVariant === "approve"
      ? "bg-emerald-600 hover:bg-emerald-500"
      : "bg-rose-600 hover:bg-rose-500";

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center p-4">
      <div
        className="absolute inset-0 bg-black/70"
        onClick={onClose}
      />
      <div className="relative w-full max-w-xl rounded-2xl border border-white/10 bg-zinc-950/80 shadow-[0_30px_120px_rgba(0,0,0,0.75)] backdrop-blur">
        <div className="flex items-center justify-between border-b border-white/10 px-5 py-4">
          <div className="text-sm font-semibold text-white">{title}</div>
          <button
            onClick={onClose}
            className="rounded-lg px-2 py-1 text-white/60 hover:bg-white/5 hover:text-white"
          >
            ✕
          </button>
        </div>

        <div className="px-5 py-4">
          <label className="mb-2 block text-sm text-white/80">เหตุผล</label>
          <textarea
            value={reason}
            onChange={(e) => setReason(e.target.value)}
            rows={5}
            className="w-full rounded-xl border border-white/10 bg-white/5 p-3 text-white placeholder:text-white/30 outline-none focus:ring-2 focus:ring-white/10"
            placeholder="กรอกเหตุผล..."
          />
        </div>

        <div className="flex justify-end gap-2 border-t border-white/10 px-5 py-4">
          <button
            onClick={onClose}
            className="rounded-xl border border-white/10 bg-white/5 px-4 py-2 text-sm text-white/80 hover:bg-white/10"
          >
            ยกเลิก
          </button>
          <button
            onClick={onConfirm}
            disabled={loading || reason.trim() === ""}
            className={`rounded-xl px-4 py-2 text-sm font-semibold text-white disabled:opacity-50 ${confirmCls}`}
          >
            {loading ? "กำลังบันทึก..." : confirmText}
          </button>
        </div>
      </div>
    </div>
  );
}
