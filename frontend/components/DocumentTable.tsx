"use client";

import type { DocStatus, DocumentItem } from "@/services/api";

function StatusPill({ status }: { status: DocStatus }) {
  const cls =
    status === "PENDING"
      ? "bg-amber-500/15 text-amber-200 ring-1 ring-amber-400/25"
      : status === "APPROVED"
      ? "bg-emerald-500/15 text-emerald-200 ring-1 ring-emerald-400/25"
      : "bg-rose-500/15 text-rose-200 ring-1 ring-rose-400/25";

  return (
    <span className={`inline-flex items-center rounded-full px-3 py-1 text-xs font-semibold ${cls}`}>
      {status}
    </span>
  );
}

export default function DocumentTable({
  tab,
  documents,
  selected,
  onToggleAll,
  onToggleOne,
}: {
  tab: DocStatus;
  documents: DocumentItem[];
  selected: Record<string, boolean>;
  onToggleAll: (checked: boolean) => void;
  onToggleOne: (id: string, checked: boolean) => void;
}) {
  const canSelect = tab === "PENDING";
  const allChecked = documents.length > 0 && documents.every((d) => !!selected[d.id]);

  return (
    <div className="max-h-[520px] overflow-auto">
      <table className="w-full table-fixed text-left text-sm">
        <thead className="sticky top-0 z-10 bg-white/8 backdrop-blur">
          <tr className="text-white/60">
            {canSelect && (
              <th className="w-12 px-6 py-3">
                <input
                  type="checkbox"
                  className="h-4 w-4 accent-white/80"
                  checked={allChecked}
                  onChange={(e) => onToggleAll(e.target.checked)}
                />
              </th>
            )}

            <th className={canSelect ? "w-40 px-6 py-3 font-medium" : "w-44 px-6 py-3 font-medium"}>
              Doc No
            </th>
            <th className="px-6 py-3 font-medium">Title</th>
            <th className="w-40 px-6 py-3 font-medium">Status</th>
            <th className="w-72 px-6 py-3 font-medium">Reason</th>
          </tr>
        </thead>

        <tbody className="divide-y divide-white/10">
          {documents.map((d) => (
            <tr key={d.id} className="hover:bg-white/5 transition">
              {canSelect && (
                <td className="px-6 py-3">
                  <input
                    type="checkbox"
                    className="h-4 w-4 accent-white/80"
                    checked={!!selected[d.id]}
                    onChange={(e) => onToggleOne(d.id, e.target.checked)}
                  />
                </td>
              )}

              <td className="px-6 py-3 font-medium text-white/85">{d.doc_no}</td>
              <td className="px-6 py-3 text-white/80 truncate">{d.title}</td>
              <td className="px-6 py-3">
                <StatusPill status={d.status} />
              </td>
              <td className="px-6 py-3 text-white/60 truncate">{d.reason?.trim() ? d.reason : "-"}</td>
            </tr>
          ))}

          {documents.length === 0 && (
            <tr>
              <td colSpan={canSelect ? 5 : 4} className="px-6 py-10 text-center text-white/50">
                ไม่พบรายการ
              </td>
            </tr>
          )}
        </tbody>
      </table>
    </div>
  );
}
