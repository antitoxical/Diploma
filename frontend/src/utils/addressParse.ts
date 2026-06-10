/** Разбор строки «Город, улица…» из профиля в поля оформления заказа. */
export function parseFullAddress(full: string): { city: string; street: string } {
  const trimmed = full.trim();
  if (!trimmed) return { city: "", street: "" };

  const commaIdx = trimmed.indexOf(",");
  if (commaIdx > 0) {
    const cityPart = trimmed.slice(0, commaIdx).trim().replace(/^г\.?\s*/iu, "");
    const streetPart = trimmed.slice(commaIdx + 1).trim();
    return { city: cityPart, street: streetPart };
  }

  return { city: "", street: trimmed };
}

export function formatFullAddress(city: string, street: string): string {
  const c = city.trim();
  const s = street.trim();
  if (c && s) return `${c}, ${s}`;
  return s || c;
}
