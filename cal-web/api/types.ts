// Temporary file hard coding expected types from APIs
// TODO: Replace with zod

export interface HomeMate {
  id: number;
  email: string;
  name: string;
  role: string;
}

export interface MateProfile {
  name: string;
  email: string;
}

export interface Home {
  name: string;
  slug: string;
  description: string;
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

export class ForbiddenError extends Error {
  public readonly statusCode: number;
  constructor(message: string) {
    super(message);
    this.statusCode = 403;
    this.name = "ForbiddenError";
    Object.setPrototypeOf(this, ForbiddenError.prototype);
  }
}
