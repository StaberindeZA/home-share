import { Row } from "@/types";

export const API_URL = process.env.API_URL || "";

const workspaceUserIdsRaw = process.env.WORKSPACE_USER_IDS || "1,2";

export const workspaceUserIds = workspaceUserIdsRaw
  .split(",")
  .map((id) => parseInt(id));

export const rows: Row[] = [
  {
    start: "09:00:00",
    end: "09:29:59",
    entryIds: [
      { id: 0, value: 0 },
      { id: 0, value: 0 },
    ],
  },
  {
    start: "09:30:00",
    end: "09:59:59",
    entryIds: [
      { id: 0, value: 0 },
      { id: 0, value: 0 },
    ],
  },
  {
    start: "10:00:00",
    end: "10:29:59",
    entryIds: [
      { id: 0, value: 0 },
      { id: 0, value: 0 },
    ],
  },
  {
    start: "10:30:00",
    end: "10:59:59",
    entryIds: [
      { id: 0, value: 0 },
      { id: 0, value: 0 },
    ],
  },
  {
    start: "11:00:00",
    end: "11:29:59",
    entryIds: [
      { id: 0, value: 0 },
      { id: 0, value: 0 },
    ],
  },
  {
    start: "11:30:00",
    end: "11:59:59",
    entryIds: [
      { id: 0, value: 0 },
      { id: 0, value: 0 },
    ],
  },
  {
    start: "12:00:00",
    end: "12:29:59",
    entryIds: [
      { id: 0, value: 0 },
      { id: 0, value: 0 },
    ],
  },
  {
    start: "12:30:00",
    end: "12:59:59",
    entryIds: [
      { id: 0, value: 0 },
      { id: 0, value: 0 },
    ],
  },
  {
    start: "13:00:00",
    end: "13:29:59",
    entryIds: [
      { id: 0, value: 0 },
      { id: 0, value: 0 },
    ],
  },
  {
    start: "13:30:00",
    end: "13:59:59",
    entryIds: [
      { id: 0, value: 0 },
      { id: 0, value: 0 },
    ],
  },
  {
    start: "14:00:00",
    end: "14:29:59",
    entryIds: [
      { id: 0, value: 0 },
      { id: 0, value: 0 },
    ],
  },
  {
    start: "14:30:00",
    end: "14:59:59",
    entryIds: [
      { id: 0, value: 0 },
      { id: 0, value: 0 },
    ],
  },
  {
    start: "15:00:00",
    end: "15:29:59",
    entryIds: [
      { id: 0, value: 0 },
      { id: 0, value: 0 },
    ],
  },
  {
    start: "15:30:00",
    end: "15:59:59",
    entryIds: [
      { id: 0, value: 0 },
      { id: 0, value: 0 },
    ],
  },
  {
    start: "16:00:00",
    end: "16:29:59",
    entryIds: [
      { id: 0, value: 0 },
      { id: 0, value: 0 },
    ],
  },
  {
    start: "16:30:00",
    end: "16:59:59",
    entryIds: [
      { id: 0, value: 0 },
      { id: 0, value: 0 },
    ],
  },
  {
    start: "17:00:00",
    end: "17:29:59",
    entryIds: [
      { id: 0, value: 0 },
      { id: 0, value: 0 },
    ],
  },
  {
    start: "17:30:00",
    end: "17:59:59",
    entryIds: [
      { id: 0, value: 0 },
      { id: 0, value: 0 },
    ],
  },
];

export const startIndexMap: Map<string, number> = new Map();

rows.forEach((row, index) => startIndexMap.set(row.start, index));
