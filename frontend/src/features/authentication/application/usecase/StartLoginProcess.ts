import { TemporaryInfoImpl } from "../../domain/entity";
import { SessionRepository } from "../../domain/repository";
import { StartLoginProcessDto } from "../dto";
import { AuthenticationService } from "../services";

export interface StartLoginProcessUsecase {
  execute(dto: StartLoginProcessDto): Promise<void>;
}

export class StartLoginProcessUsecaseImpl implements StartLoginProcessUsecase {
  constructor(
    private authenticationService: AuthenticationService,
    private sessionRepository: SessionRepository
  ) {}

  async execute(dto: StartLoginProcessDto): Promise<void> {
    const { provider, rememberMe, redirectToAuthenticationPage } = dto;

    const loginUrlResponse = await this.authenticationService.getLoginUrl({
      providerName: provider,
    });

    const temporaryInfo = TemporaryInfoImpl.new(
      provider,
      loginUrlResponse.state,
      rememberMe
    );

    // Save temporary info to session storage
    await this.sessionRepository.saveTemporaryInfo(temporaryInfo);

    // Redirect to external authentication page
    redirectToAuthenticationPage(loginUrlResponse.authUrl);
  }
}
