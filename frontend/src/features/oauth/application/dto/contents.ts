export interface Content {
  get name(): string;
  get path(): string;
  get type(): string;
  get size(): number;
  get url(): string;
}
