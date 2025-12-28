import { useParams, useSearchParams } from "react-router-dom";
import { RouteError } from "@/ui/error/RouteError";
import { Page } from "../types/Page";

type ParamConfig<T extends string | never> = {
  name: T;
  required?: boolean;
  validate?: (value: string) => boolean;
};

interface TypedParamsGuardProps<
  TPath extends string | never,
  TQuery extends string | never
> {
  pathParams: TPath[] extends never[] ? undefined : ParamConfig<TPath>[];
  queryParams: TQuery[] extends never[] ? undefined : ParamConfig<TQuery>[];
  fallbackPath: string;
  errorMessage?: string;
  element: Page<TPath, TQuery>;
}

export function TypedParamsGuard<
  TPath extends string | never = never,
  TQuery extends string | never = never
>({
  pathParams,
  queryParams,
  fallbackPath,
  errorMessage,
  element,
}: TypedParamsGuardProps<TPath, TQuery>) {
  const allPathParams = useParams();
  const [searchParams] = useSearchParams();

  const errors: string[] = [];
  const pathResult: Record<string, string> = {};
  const queryResult: Record<string, string> = {};

  if (typeof pathParams !== "undefined") {
    // パスパラメータ処理
    for (const config of pathParams) {
      const value = allPathParams[config.name];

      if (!value && config.required !== false) {
        errors.push(`Missing path parameter: ${config.name}`);
        continue;
      }

      if (value) {
        if (config.validate && !config.validate(value)) {
          errors.push(`Invalid path parameter: ${config.name}`);
          continue;
        }
        pathResult[config.name] = value;
      }
    }
  }

  if (typeof queryParams !== "undefined") {
    // クエリパラメータ処理
    for (const config of queryParams) {
      const value = searchParams.get(config.name);

      if (!value && config.required !== false) {
        errors.push(`Missing query parameter: ${config.name}`);
        continue;
      }

      if (value) {
        if (config.validate && !config.validate(value)) {
          errors.push(`Invalid query parameter: ${config.name}`);
          continue;
        }
        queryResult[config.name] = value;
      }
    }
  }

  if (errors.length > 0) {
    return (
      <RouteError
        message={errorMessage || errors.join(", ")}
        fallbackPath={fallbackPath}
      />
    );
  }

  // Type-safe assertion since we've validated all params
  const validatedPath = pathResult as TPath extends never
    ? never
    : Record<TPath, string>;
  const validatedQuery = queryResult as TQuery extends never
    ? never
    : Record<TQuery, string>;

  return element({
    path: validatedPath,
    query: validatedQuery,
  });
}
