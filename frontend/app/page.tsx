"use client";

import { useEffect, useMemo, useState } from "react";
import ApprovalModal from "@/components/ApprovalModal";
import DocumentTable from "@/components/DocumentTable";
import {
  approveDocuments,
  listDocuments,
  rejectDocuments,
  type DocStatus,
  type DocumentItem,
} from "@/services/api";

type Tab = DocStatus;

function cn(...xs: Array<string | false | undefined>) {
  return xs.filter(Boolean).join(" ");
}

export default function Page() {
  const [tab, setTab] = useState<Tab>("PENDING");
  const [loading, setLoading] = useState(false);
  const [docs, setDocs] = useState<DocumentItem[]>([]);
  const [selected, setSelected] = useState<Record<string, boolean>>({});

  // modal
  const [modalOpen, setModalOpen] = useState(false);
  const [modalType, setModalType] = useState<"approve" | "reject">("approve");
  const [reason, setReason] = useState("");
  const [saving, setSaving] = useState(false);

  const filtered = useMemo(() => docs.filter((d) => d.status === tab), [docs, tab]);
  const selectedIds = useMemo(
    () => Object.entries(selected).filter(([, v]) => v).map(([k]) => k),
    [selected]
  );

  const canSelect = tab === "PENDING";
  const canBulkAction = canSelect && selectedIds.length > 0;

  async function refresh() {
    setLoading(true);
    try {
      const res = await listDocuments(); // { statusCode, data }
      setDocs(res.data || []);
    } finally {
      setLoading(false);
    }
  }

  useEffect(() => {
    refresh();
  }, []);

  function clearSelection() {
    setSelected({});
  }

  function openApprove() {
    setModalType("approve");
    setReason("");
    setModalOpen(true);
  }

  function openReject() {
    setModalType("reject");
    setReason("");
    setModalOpen(true);
  }

  async function onConfirm() {
    if (!canBulkAction) return;
    const r = reason.trim();
    if (!r) return;

    setSaving(true);
    try {
      if (modalType === "approve") {
        await approveDocuments(selectedIds, r);
      } else {
        await rejectDocuments(selectedIds, r);
      }

      setModalOpen(false);
      clearSelection();
      await refresh();
    } catch (e: any) {
      alert(e?.message ?? "Something went wrong");
    } finally {
      setSaving(false);
    }
  }

  return (
    <main className="min-h-screen it03-bg text-white">
      <div className="mx-auto max-w-6xl px-6 py-10">
        {/* Header + Actions */}
        <div className="flex items-start justify-between gap-6">
          <div>
            <h1 className="text-3xl font-semibold tracking-tight">Workflow Approval</h1>
            <p className="mt-1 text-sm text-white/60">จัดการอนุมัติ/ไม่อนุมัติ</p>

            {/* Tabs */}
            <div className="mt-5 inline-flex rounded-xl bg-white/5 p-1 ring-1 ring-white/10">
              <TabButton
                active={tab === "PENDING"}
                onClick={() => {
                  setTab("PENDING");
                  clearSelection();
                  setModalOpen(false);
                }}
              >
                รออนุมัติ
              </TabButton>
              <TabButton
                active={tab === "APPROVED"}
                onClick={() => {
                  setTab("APPROVED");
                  clearSelection();
                  setModalOpen(false);
                }}
              >
                อนุมัติ
              </TabButton>
              <TabButton
                active={tab === "REJECTED"}
                onClick={() => {
                  setTab("REJECTED");
                  clearSelection();
                  setModalOpen(false);
                }}
              >
                ไม่อนุมัติ
              </TabButton>
            </div>
          </div>

          <div className="flex items-center gap-2">
            <button
              onClick={refresh}
              className="rounded-xl bg-white/5 px-4 py-2 text-sm ring-1 ring-white/10 hover:bg-white/10"
            >
              {loading ? "กำลังโหลด..." : "รีเฟรช"}
            </button>

            <button
              onClick={clearSelection}
              className="rounded-xl bg-white/5 px-4 py-2 text-sm ring-1 ring-white/10 hover:bg-white/10"
            >
              เคลียร์
            </button>

            <button
              disabled={!canBulkAction}
              onClick={openApprove}
              className={cn(
                "ml-2 rounded-xl px-4 py-2 text-sm font-semibold ring-1 transition",
                canBulkAction
                  ? "bg-emerald-500/15 text-emerald-200 ring-emerald-400/25 hover:bg-emerald-500/20"
                  : "bg-emerald-500/10 text-emerald-200/30 ring-emerald-400/10 cursor-not-allowed"
              )}
            >
              อนุมัติ
            </button>

            <button
              disabled={!canBulkAction}
              onClick={openReject}
              className={cn(
                "rounded-xl px-4 py-2 text-sm font-semibold ring-1 transition",
                canBulkAction
                  ? "bg-rose-500/15 text-rose-200 ring-rose-400/25 hover:bg-rose-500/20"
                  : "bg-rose-500/10 text-rose-200/30 ring-rose-400/10 cursor-not-allowed"
              )}
            >
              ไม่อนุมัติ
            </button>
          </div>
        </div>

        <div className="mt-8 glass rounded-2xl overflow-hidden">
          <div className="flex items-center justify-between px-6 py-4">
            <div className="text-sm text-white/70">
              STATUS : <span className="text-white">{tab}</span>
              {!canSelect && <span className="ml-2 text-white/40">(เลือกไม่ได้)</span>}
            </div>

            <div className="text-sm text-white/60">
              ทั้งหมด: <span className="text-white/80">{filtered.length}</span> รายการ
            </div>
          </div>

          <div className="border-t border-white/10" />

          <DocumentTable
            tab={tab}
            documents={filtered}
            selected={selected}
            onToggleAll={(checked) => {
              if (!canSelect) return;
              const next: Record<string, boolean> = {};
              for (const d of filtered) next[d.id] = checked;
              setSelected(next);
            }}
            onToggleOne={(id, checked) => {
              if (!canSelect) return;
              setSelected((prev) => ({ ...prev, [id]: checked }));
            }}
          />

          {canSelect && (
            <div className="px-6 py-3 text-xs text-white/50 border-t border-white/10">
              เลือกแล้ว: <span className="text-white/80">{selectedIds.length}</span> รายการ
            </div>
          )}
        </div>
      </div>

      <ApprovalModal
        open={modalOpen}
        title={modalType === "approve" ? "ยืนยันการอนุมัติ" : "ยืนยันการไม่อนุมัติ"}
        confirmText={modalType === "approve" ? "อนุมัติ" : "ไม่อนุมัติ"}
        confirmVariant={modalType}
        reason={reason}
        setReason={setReason}
        onClose={() => setModalOpen(false)}
        onConfirm={onConfirm}
        loading={saving}
      />
    </main>
  );
}

function TabButton({
  active,
  children,
  onClick,
}: {
  active: boolean;
  children: React.ReactNode;
  onClick: () => void;
}) {
  return (
    <button
      onClick={onClick}
      className={cn(
        "rounded-lg px-4 py-2 text-sm transition",
        active ? "bg-white text-black shadow font-semibold" : "text-white/70 hover:text-white hover:bg-white/10"
      )}
    >
      {children}
    </button>
  );
}
