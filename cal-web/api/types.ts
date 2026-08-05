// Temporary file hard coding expected types from APIs

export interface HomeMate {
  id: string;
  email: string;
  name: string;
}

export class UnauthorizedError extends Error {
  public readonly statusCode: number;
  constructor(message: string) {
    super(message);
    this.statusCode = 401;
    this.name = "UnauthorizedError";
    Object.setPrototypeOf(this, UnauthorizedError.prototype);
  }
}
