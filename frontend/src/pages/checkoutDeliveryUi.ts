import type { DeliveryMethod } from "../api/types";

export type DeliveryUiKind = "courier" | "pickup" | "pvz" | "postal";

const VALID_KINDS: DeliveryUiKind[] = ["courier", "pickup", "pvz", "postal"];

function isDeliveryUiKind(v: string): v is DeliveryUiKind {
  return (VALID_KINDS as string[]).includes(v);
}

/** Предпочитает kind из API; иначе угадывает по названию (для старых записей). */
export function resolveDeliveryKind(method: DeliveryMethod | undefined): DeliveryUiKind {
  if (method?.kind && isDeliveryUiKind(method.kind)) {
    return method.kind;
  }
  return classifyDeliveryUi(method);
}

/** Запасной вариант, если kind в БД не задан. */
export function classifyDeliveryUi(method: DeliveryMethod | undefined): DeliveryUiKind {
  if (!method) return "courier";
  const hay = `${method.name}\n${method.description ?? ""}`.toLowerCase();

  if (/\bсамовывоз\b/.test(hay) || /\bpickup\b/i.test(method.name)) return "pickup";

  // ПВЗ в описании самовывоза («уточните адрес ПВЗ у менеджера») не переводит в СДЭК.
  if (
    /\bсдэк\b/.test(hay) ||
    /\bcdek\b/i.test(hay) ||
    (/\bпвз\b/.test(hay) && !/\bсамовывоз\b/.test(hay))
  )
    return "pvz";

  if (/\bпочт/i.test(hay) || /\bпочта\b/.test(hay)) return "postal";

  return "courier";
}

export function needsStreetAddress(kind: DeliveryUiKind): boolean {
  return kind === "courier" || kind === "postal";
}

export function cityFieldLabel(kind: DeliveryUiKind): string {
  switch (kind) {
    case "pvz":
      return "Город пункта выдачи*";
    case "pickup":
      return "Ваш город*";
    default:
      return "Город доставки*";
  }
}

export function cityFieldPlaceholder(kind: DeliveryUiKind): string {
  switch (kind) {
    case "pvz":
      return "Город, где находится ПВЗ";
    case "pickup":
      return "Для связи с менеджером";
    default:
      return "Например, Смоленск";
  }
}

export function buildAddressSnapshot(params: {
  kind: DeliveryUiKind;
  methodLabel: string;
  city: string;
  streetAddress: string;
  comment: string;
  pickupPointLine: string;
}): string {
  const lines: string[] = [];
  const city = params.city.trim();
  const street = params.streetAddress.trim();
  const comment = params.comment.trim();
  const pvz = params.pickupPointLine.trim();

  lines.push(`Способ доставки: ${params.methodLabel}`);
  if (city) lines.push(`Город: ${city}`);

  if (params.kind === "pickup") {
    lines.push("Тип: самовывоз");
  } else if (params.kind === "pvz") {
    lines.push("Тип: пункт выдачи (ПВЗ)");
    if (pvz) lines.push(`ПВЗ: ${pvz}`);
  } else {
    if (street) lines.push(`Адрес: ${street}`);
  }

  if (comment) lines.push(`Комментарий: ${comment}`);
  return lines.join("\n");
}
