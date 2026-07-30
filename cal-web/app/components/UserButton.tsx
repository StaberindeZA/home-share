'use client'

import { useState } from "react";
import { addEntry, changeEntry } from "../actions";
import { RowEntry, User } from "../types";

interface UserButtonProps {
  startTime: string,
  endTime: string,
  user: User,
  rowEntry?: RowEntry
}

const valueMapper = ["open", "BOOKED", "Talking", "Not Talking"]

export default function UserButton({ startTime, endTime, user, rowEntry }: UserButtonProps) {
  const entryId = rowEntry?.id ?? 0;
  const entryValue = rowEntry?.value ?? 0;

  const [loading, setLoading] = useState(false);
  const onClickHandler = async () => {
    try {
      if (!entryValue && entryValue != 0) {
        setLoading(true)
        await addEntry(user.id, startTime, endTime)
      } else if (entryValue >= 0 && entryValue < 3 && entryId) {
        setLoading(true)
        await changeEntry(entryId, entryValue + 1)
      } else if (entryValue === 3 && entryId) {
        setLoading(true)
        await changeEntry(entryId, 0)
      } else {
        throw new Error(`Unexpected values. Value: ${entryValue}, ID: ${entryId}`)
      }
    } catch (error) {
      console.error(error);
      throw error
    } finally {
      setLoading(false)
    }
  }

  const entryValueDisplay = entryValue ? valueMapper[entryValue] : valueMapper[0];

  return (
    <div
      className="grid h-8 items-center justify-center active:bg-gray-700 select-none"
      onClick={onClickHandler}>{!loading ? entryValueDisplay : "Loading..."}</div>
  )
}

