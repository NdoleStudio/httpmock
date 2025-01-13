import { AxiosError } from "axios";
import Bag from "@/utils/bag";
import { capitalize } from "@/utils/filters";

export class ErrorMessages extends Bag<string> {}

const sanitize = (key: string, values: Array<string>): Array<string> => {
  return values.map((value: string) => {
    return capitalize(
      value
        .split(key)
        .join(key.replace("_", " "))
        .split("_")
        .join(" ")
        .split("-")
        .join(" ")
        .split(" char")
        .join(" character")
        .split(" field ")
        .join(" "),
    );
  });
};

export const getErrorMessages = (error: AxiosError): ErrorMessages => {
  const errors = new ErrorMessages();
  if (
    error === null ||
    /* eslint-disable @typescript-eslint/no-explicit-any */
    typeof (error.response?.data as any)?.data !== "object" ||
    /* eslint-disable @typescript-eslint/no-explicit-any */
    (error.response?.data as any)?.data === null ||
    error.response?.status !== 422
  ) {
    return errors;
  }

  /* eslint-disable @typescript-eslint/no-explicit-any */
  Object.keys((error.response?.data as any).data).forEach((key: string) => {
    /* eslint-disable @typescript-eslint/no-explicit-any */
    errors.addMany(key, sanitize(key, (error.response?.data as any).data[key]));
  });

  return errors;
};
