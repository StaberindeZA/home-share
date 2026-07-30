'use server'

import { revalidatePath } from "next/cache";
import { createEntry, deleteEntry, fetchListEntries, updateEntry } from "../api";
import { rows, startIndexMap, users } from "../constants";
import { convertTimeToUTCISO, convertToLocalTime, wait } from "../utils";
import { redirect } from "next/navigation";

export async function addEntry(userId: number, startTime: string, endTime: string) {
	const startDateString = convertTimeToUTCISO(startTime);
	const endDateString = convertTimeToUTCISO(endTime);
	await createEntry(userId, startDateString, endDateString)

	revalidatePath('/')
}

export async function changeEntry(entryId: number, entryValue: number) {
	await updateEntry(entryId, entryValue);

	revalidatePath('/')
}

export async function removeEntry(entryId: number) {
	await deleteEntry(entryId);

	revalidatePath('/')
}

export async function listEntries(timeZone: string) {
	const currentRows = structuredClone(rows)
	const entries = await Promise.all([
		fetchListEntries(users[0].id),
		fetchListEntries(users[1].id)
	])


	entries.forEach((entry, userIndex) => {
		entry.forEach((e: any) => {
			const startTime = convertToLocalTime(new Date(e.start), timeZone);
			const startIndex = startIndexMap.get(startTime);

			if (!!startIndex || startIndex === 0) {
				currentRows[startIndex].entryIds[userIndex] = { id: e.id, value: e.value }
			} else {
				console.error("Could not find index for:", startTime)
			}

		})
	})

	return currentRows;
}

export async function saveOptions(formData: FormData) {
	const checkedUser = formData.get('user-checkbox')

	if (checkedUser === null) {
		redirect(`/?user=${users[0].id}`)
	} else if (checkedUser === 'on') {
		redirect(`/?user=${users[1].id}`)
	} else {
		throw new Error('Unexpected checkbox value');
	}
}
