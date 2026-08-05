export type RowEntry = {
  id: number;
  value: number;
};

export type Row = {
  start: string;
  end: string;
  entryIds: RowEntry[];
};

export type SearchParams = { [key: string]: string | string[] | undefined };
export type SearchParamsPromise = Promise<SearchParams>;
