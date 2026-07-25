export type User = {
	id: number;
	name: string;
}

export type Row = {
	start: string;
	end: string;
	entryIds: number[];
}

export type SearchParams = { [key: string]: string | string[] | undefined };
export type SearchParamsPromise = Promise<SearchParams>;


