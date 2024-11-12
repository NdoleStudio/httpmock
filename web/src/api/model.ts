/* eslint-disable */
/* tslint:disable */
/*
 * ---------------------------------------------------------------
 * ## THIS FILE WAS GENERATED VIA SWAGGER-TYPESCRIPT-API        ##
 * ##                                                           ##
 * ## AUTHOR: acacode                                           ##
 * ## SOURCE: https://github.com/acacode/swagger-typescript-api ##
 * ---------------------------------------------------------------
 */

export interface EntitiesProject {
  /** @example "2022-06-05T14:26:02.302718+03:00" */
  created_at: string;
  /** @example "Mock API for an online store for selling shoes" */
  description: string;
  /** @example "8f9c71b8-b84e-4417-8408-a62274f65a08" */
  id: string;
  /** @example "Mock Stripe API" */
  name: string;
  /** @example "api" */
  subdomain: string;
  /** @example "2022-06-05T14:26:10.303278+03:00" */
  updated_at: string;
  /** @example "WB7DRDWrJZRGbYrv2CKGkqbzvqdC" */
  user_id: string;
}

export interface RequestsProjectCreateRequest {
  description: string;
  name: string;
  subdomain: string;
}

export interface RequestsProjectUpdateRequest {
  description: string;
  name: string;
}

export interface ResponsesBadRequest {
  /** @example "The request body is not a valid JSON string" */
  data: string;
  /** @example "The request isn't properly formed" */
  message: string;
  /** @example "error" */
  status: string;
}

export interface ResponsesInternalServerError {
  /** @example "We ran into an internal error while handling the request." */
  message: string;
  /** @example "error" */
  status: string;
}

export interface ResponsesNoContent {
  /** @example "item deleted successfully" */
  message: string;
  /** @example "success" */
  status: string;
}

export interface ResponsesNotFound {
  /** @example "cannot find item with ID [32343a19-da5e-4b1b-a767-3298a73703ca]" */
  message: string;
  /** @example "error" */
  status: string;
}

export interface ResponsesOkArrayEntitiesProject {
  data: EntitiesProject[];
  /** @example "Request handled successfully" */
  message: string;
  /** @example "success" */
  status: string;
}

export interface ResponsesOkEntitiesProject {
  data: EntitiesProject;
  /** @example "Request handled successfully" */
  message: string;
  /** @example "success" */
  status: string;
}

export interface ResponsesUnauthorized {
  /** @example "Make sure your Bearer token is set in the [Bearer] header in the request" */
  data: string;
  /** @example "You are not authorized to carry out this request." */
  message: string;
  /** @example "error" */
  status: string;
}

export interface ResponsesUnprocessableEntity {
  data: Record<string, string[]>;
  /** @example "validation errors while sending message" */
  message: string;
  /** @example "error" */
  status: string;
}
