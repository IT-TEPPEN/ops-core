import { SessionRepository } from "../../domain/repository";
import { AuthenticationService } from "../services";

export interface SignOutUsecase {
  execute(): Promise<void>;
}

export class SignOutUsecaseImpl implements SignOutUsecase {
  constructor(
    private sessionRepository: SessionRepository,
    private authenticationService: AuthenticationService
  ) {}

  async execute(): Promise<void> {
    await this.authenticationService.logout();

    await this.sessionRepository.removeToken();
  }
}
