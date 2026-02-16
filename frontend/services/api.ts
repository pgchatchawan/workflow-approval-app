const API_BASE = process.env.NEXT_PUBLIC_API_BASE ?? "http://localhost:8080";

export type DocStatus = "PENDING" | "APPROVED" | "REJECTED";

export type DocumentItem = {
  id: string;
  doc_no: string;
  title: string;
  status: DocStatus;
  reason?: string;
  created_at: string;
  updated_at: string;
};

export async function listDocuments(status?: DocStatus) {
  const url = status
    ? `${API_BASE}/api/documents?status=${status}`
    : `${API_BASE}/api/documents`;

  const res = await fetch(url, { cache: "no-store" });
  if (!res.ok) throw new Error("Failed to fetch documents");
  return (await res.json()) as { statusCode: number; data: DocumentItem[] };
}

export async function approveDocuments(document_ids: string[], reason: string) {
  const res = await fetch(`${API_BASE}/api/documents/approval`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ document_ids, reason }),
  });
  if (!res.ok) throw new Error("Failed to approve");
  return (await res.json()) as {
    statusCode: number;
    message: string;
    requested: number;
    approved: number;
  };
}

export async function rejectDocuments(document_ids: string[], reason: string) {
  const res = await fetch(`${API_BASE}/api/documents/rejection`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ document_ids, reason }),
  });
  if (!res.ok) throw new Error("Failed to reject");
  return (await res.json()) as {
    statusCode: number;
    message: string;
    requested: number;
    rejected: number;
  };
}
