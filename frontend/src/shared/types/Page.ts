export type Page<
  TPath extends string | never = never,
  TQuery extends string | never = never
> = (params: {
  path: [TPath] extends [never] ? never : Record<TPath, string>;
  query: [TQuery] extends [never] ? never : Record<TQuery, string>;
}) => React.ReactNode;
