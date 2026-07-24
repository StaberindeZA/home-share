'use client'

import { addEntry, removeEntry } from "../actions";
import { User } from "../types";

interface UserButtonProps {
  startTime: string,
  endTime: string,
  user: User,
  entryId?: number,
}

export default function UserButton(props: UserButtonProps) {
  return (
    <div onClick={async () => {
      props.entryId ? removeEntry(props.entryId) : addEntry(props.user.id, props.startTime, props.endTime)
    }}>{props.entryId ? "BOOKED" : "open"}</div>
  )
}

