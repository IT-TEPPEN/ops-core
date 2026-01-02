export interface DocumentVariableBase<T extends string> {
  name: string;
  label: string;
  description: string | null;
  required: boolean;
  type: T;
}

export interface DocumentVariableString extends DocumentVariableBase<"string"> {
  defaultValue: string | null;
}

export interface DocumentVariableNumber extends DocumentVariableBase<"number"> {
  defaultValue: number | null;
}

export interface DocumentVariableBoolean
  extends DocumentVariableBase<"boolean"> {
  defaultValue: boolean | null;
}

export interface DocumentVariableDate extends DocumentVariableBase<"date"> {
  defaultValue: string | null; // ISO date string
}

export type DocumentVariable =
  | DocumentVariableString
  | DocumentVariableNumber
  | DocumentVariableBoolean
  | DocumentVariableDate;

export interface DocumentProcedureMeta {
  title: string;
  owner: string;
  type: "procedure";
  tags: string[];
  variables: DocumentVariable[];
}

export interface DocumentKnowledgeMeta {
  title: string;
  owner: string;
  type: "knowledge";
  tags: string[];
}
