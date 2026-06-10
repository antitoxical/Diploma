export type FieldErrors = Record<string, string>;

const NAME_MAX = 50;
const ADDRESS_MIN = 10;
const ADDRESS_MAX = 300;
const CITY_MAX = 100;
const MAX_ADDRESSES = 20;

const namePattern = /^[\p{L}\s\-']+$/u;
const cityPattern = /^[\p{L}\s\-.]+$/u;

const personNameInputFilter = /[^\p{L}\s\-']/gu;
const cityInputFilter = /[^\p{L}\s\-.]/gu;

export function filterPersonNameInput(value: string): string {
  return value.replace(personNameInputFilter, "");
}

export function filterCityInput(value: string): string {
  return value.replace(cityInputFilter, "");
}

function invalidCharsMessage(label: string, allowed: string): string {
  return `${label}: только ${allowed}, без цифр и спецсимволов`;
}

export function normalizePhone(value: string): string {
  return value.replace(/\D/g, "");
}

const PHONE_MIN_DIGITS = 10;
const PHONE_MAX_DIGITS = 15;

export function validatePhone(value: string, required = true): string | null {
  const trimmed = value.trim();
  if (!trimmed) return required ? "Укажите номер телефона" : null;

  const digits = normalizePhone(trimmed);

  if (digits.length < PHONE_MIN_DIGITS || digits.length > PHONE_MAX_DIGITS) {
    return `Введите номер с кодом страны: от ${PHONE_MIN_DIGITS} до ${PHONE_MAX_DIGITS} цифр, например +375 29 123-45-67`;
  }

  return null;
}

export function validatePersonName(
  value: string,
  label: string,
  required = true,
): string | null {
  const trimmed = value.trim();
  if (!trimmed) return required ? `Укажите ${label.toLowerCase()}` : null;
  if (trimmed.length < 2) return `${label}: минимум 2 символа`;
  if (trimmed.length > NAME_MAX) return `${label}: не более ${NAME_MAX} символов`;
  if (/\d/.test(trimmed) || !namePattern.test(trimmed)) {
    return invalidCharsMessage(label, "буквы, пробелы, дефис");
  }
  return null;
}

export function validateFullName(value: string): string | null {
  const trimmed = value.trim();
  const base = validatePersonName(trimmed, "Ф.И.О.", true);
  if (base) return base;
  const parts = trimmed.split(/\s+/).filter(Boolean);
  if (parts.length < 2) return "Укажите фамилию и имя (минимум два слова)";
  return null;
}

export function validateAddress(value: string, required = true): string | null {
  const trimmed = value.trim();
  if (!trimmed) return required ? "Укажите адрес" : null;
  if (trimmed.length < ADDRESS_MIN) {
    return `Адрес слишком короткий (минимум ${ADDRESS_MIN} символов)`;
  }
  if (trimmed.length > ADDRESS_MAX) {
    return `Адрес слишком длинный (максимум ${ADDRESS_MAX} символов)`;
  }
  return null;
}

export function validateCity(value: string): string | null {
  const trimmed = value.trim();
  if (!trimmed) return "Укажите город";
  if (trimmed.length < 2) return "Название города слишком короткое";
  if (trimmed.length > CITY_MAX) return `Название города слишком длинное (максимум ${CITY_MAX} символов)`;
  if (/\d/.test(trimmed) || !cityPattern.test(trimmed)) {
    return invalidCharsMessage("Город", "буквы, пробелы, дефис и точка");
  }
  return null;
}

export function validateStreetWithoutCity(street: string, city: string): string | null {
  const s = street.trim().toLowerCase();
  const c = city.trim().toLowerCase();
  if (!s || !c) return null;
  if (s === c) {
    return "Укажите улицу и дом — город уже указан в поле «Город»";
  }
  if (s.startsWith(c + ",") || s.startsWith(c + " ")) {
    return "Не дублируйте город — укажите только улицу, дом, подъезд";
  }
  return null;
}

export function validateProfileForm(fields: {
  firstName: string;
  lastName: string;
  phone: string;
}): FieldErrors {
  const errors: FieldErrors = {};
  const first = validatePersonName(fields.firstName, "Имя", true);
  if (first) errors.firstName = first;
  const last = validatePersonName(fields.lastName, "Фамилия", false);
  if (last) errors.lastName = last;
  const phone = validatePhone(fields.phone, true);
  if (phone) errors.phone = phone;
  return errors;
}

export function validateNewAddress(value: string, currentCount: number): string | null {
  if (currentCount >= MAX_ADDRESSES) {
    return `Можно сохранить не более ${MAX_ADDRESSES} адресов`;
  }
  return validateAddress(value, true);
}

export function validateCheckoutForm(fields: {
  name: string;
  phone: string;
  city: string;
  streetAddress: string;
  pickupPointLine: string;
  courierLike: boolean;
  deliveryKind: string;
  selectedDelivery: string;
  selectedPayment: string;
  cartEmpty: boolean;
}): FieldErrors {
  const errors: FieldErrors = {};

  const nameErr = validateFullName(fields.name);
  if (nameErr) errors.name = nameErr;

  const phoneErr = validatePhone(fields.phone, true);
  if (phoneErr) errors.phone = phoneErr;

  const cityErr = validateCity(fields.city);
  if (cityErr) errors.city = cityErr;

  if (fields.courierLike) {
    const addrErr = validateAddress(fields.streetAddress, true);
    if (addrErr) errors.streetAddress = addrErr;
    const dupErr = validateStreetWithoutCity(fields.streetAddress, fields.city);
    if (dupErr) errors.streetAddress = dupErr;
  }

  if (fields.deliveryKind === "pvz") {
    const pvzErr = validateAddress(fields.pickupPointLine, true);
    if (pvzErr) errors.pickupPointLine = pvzErr;
  }

  if (!fields.selectedDelivery || !fields.selectedPayment) {
    errors.delivery = "Выберите способ доставки и оплаты";
  }

  if (fields.cartEmpty) {
    errors.cart = "Корзина пуста";
  }

  return errors;
}

export function firstErrorMessage(errors: FieldErrors): string | null {
  const keys = Object.keys(errors);
  return keys.length > 0 ? errors[keys[0]] : null;
}
