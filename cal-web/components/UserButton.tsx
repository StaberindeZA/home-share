"use client";

import { useState } from "react";
import { addEntry, changeEntry } from "@/actions";
import { RowEntry } from "@/types";
import { User } from "next-auth";
import { HomeMate } from "@/api/types";

interface UserButtonProps {
  startTime: string;
  endTime: string;
  homeMate: HomeMate;
  sessionUser?: User;
  rowEntry?: RowEntry;
}

const valueMapper = ["open", "BOOKED", "Talking", "Not Talking"];

function getButtonColor(entryValue: number) {
  switch (entryValue) {
    case 1:
      return "bg-amber-300 dark:bg-amber-700";
    case 2:
      return "bg-red-300 dark:bg-red-700";
    case 3:
      return "bg-sky-300 dark:bg-sky-700";
    case 0:
    default:
      return "color-bg-light dark:color-bg-dark";
  }
}

export default function UserButton({
  startTime,
  endTime,
  homeMate,
  sessionUser,
  rowEntry,
}: UserButtonProps) {
  const entryId = rowEntry?.id ?? 0;
  const entryValue = rowEntry?.value ?? 0;

  const [loading, setLoading] = useState(false);
  const onClickHandler = async () => {
    if (homeMate.email != sessionUser?.email) return;
    if (loading) return;

    try {
      if (!entryId) {
        setLoading(true);
        await addEntry(startTime, endTime);
      } else if (entryValue >= 0 && entryValue < 3 && entryId) {
        setLoading(true);
        await changeEntry(entryId, entryValue + 1);
      } else if (entryValue === 3 && entryId) {
        setLoading(true);
        await changeEntry(entryId, 0);
      } else {
        throw new Error(
          `Unexpected values. Value: ${entryValue}, ID: ${entryId}`,
        );
      }
    } catch (error) {
      console.error(error);
      throw error;
    } finally {
      setLoading(false);
    }
  };

  const entryValueDisplay = entryValue
    ? valueMapper[entryValue]
    : valueMapper[0];

  const buttonColor = getButtonColor(entryValue);

  return (
    <div
      className={`grid h-8 items-center justify-center active:bg-gray-200 active:dark:bg-gray-700 select-none ${buttonColor}`}
      onClick={onClickHandler}
    >
      {entryValueDisplay}
    </div>
  );
}
