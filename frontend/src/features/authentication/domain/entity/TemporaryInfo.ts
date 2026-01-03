export interface TemporaryInfo {
  get provider(): string;
  get state(): string;
  get rememberMe(): boolean;
  toJSONString(): string;
}

export class TemporaryInfoImpl implements TemporaryInfo {
  private constructor(
    public provider: string,
    public state: string,
    public rememberMe: boolean
  ) {}

  static new(
    provider: string,
    state: string,
    rememberMe: boolean
  ): TemporaryInfo {
    return new TemporaryInfoImpl(provider, state, rememberMe);
  }

  static reconstruct(
    provider: string,
    state: string,
    rememberMe: boolean
  ): TemporaryInfo {
    return new TemporaryInfoImpl(provider, state, rememberMe);
  }

  static fromJSONString(jsonString: string): TemporaryInfo {
    const obj = JSON.parse(jsonString);
    return new TemporaryInfoImpl(obj.provider, obj.state, obj.rememberMe);
  }

  toJSONString(): string {
    return JSON.stringify({
      provider: this.provider,
      state: this.state,
      rememberMe: this.rememberMe,
    });
  }
}
