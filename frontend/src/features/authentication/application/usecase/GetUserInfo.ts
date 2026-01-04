import { IdentityDto } from "../dto";
import { AuthenticationService } from "../services";

export interface GetUserIdentityUsecase {
  execute(): Promise<IdentityDto | null>;
}

export class GetUserIdentityUsecaseImpl implements GetUserIdentityUsecase {
  constructor(private authenticationService: AuthenticationService) {}

  execute(): Promise<IdentityDto | null> {
    return this.authenticationService.getIdentity();
  }
}
