import { Pagenation } from "../../../shared/data-types";
import { Repository } from "./Repository";

export interface Repositories {
  data: Repository[];
  pagenation: Pagenation;
}
