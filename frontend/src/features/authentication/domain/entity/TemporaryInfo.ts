export interface TemporaryInfo {
  get provider(): string;
  get state(): string;
  get rememberMe(): boolean;
  get from(): string;
  toJSONString(): string;
}

export class TemporaryInfoImpl implements TemporaryInfo {
  private constructor(
    public readonly provider: string,
    public readonly state: string,
    public readonly rememberMe: boolean,
    public readonly from: string
  ) {}

  static new(
    provider: string,
    state: string,
    rememberMe: boolean,
    from: string
  ): TemporaryInfo {
    return new TemporaryInfoImpl(provider, state, rememberMe, from);
  }

  static reconstruct(
    provider: string,
    state: string,
    rememberMe: boolean,
    from: string
  ): TemporaryInfo {
    return new TemporaryInfoImpl(provider, state, rememberMe, from);
  }

  static fromJSONString(jsonString: string): TemporaryInfo {
    const obj = JSON.parse(jsonString);
    return new TemporaryInfoImpl(
      obj.provider,
      obj.state,
      obj.rememberMe,
      obj.from
    );
  }

  toJSONString(): string {
    return JSON.stringify({
      provider: this.provider,
      state: this.state,
      rememberMe: this.rememberMe,
      from: this.from,
    });
  }
}
