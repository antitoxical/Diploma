const ERROR_MESSAGES: Record<string, string> = {
  error_create_order: "Не удалось оформить заказ",
  error_invalid_request_body: "Проверьте введённые данные",
};

export function formatApiError(statusCode: number, bodyText: string): string {
  const raw = bodyText.trim();
  let status: string | undefined;
  let data: unknown;

  try {
    const json = JSON.parse(raw);
    status = json.status;
    data = json.data;
  } catch {
    if (raw.length > 0 && raw.length < 160) return raw;
    return defaultMessage(statusCode);
  }

  const detail = typeof data === "string" ? data : "";

  if (status === "error_create_order" && /insufficient stock/i.test(detail)) {
    return "Недостаточно товара на складе. Уменьшите количество в корзине.";
  }

  if (detail && detail.length < 160 && !detail.includes("{")) {
    if (/insufficient stock/i.test(detail)) {
      return "Недостаточно товара на складе. Уменьшите количество в корзине.";
    }
    return detail;
  }

  if (status && ERROR_MESSAGES[status]) {
    return ERROR_MESSAGES[status];
  }

  return defaultMessage(statusCode);
}

function defaultMessage(statusCode: number): string {
  if (statusCode === 400) return "Некорректный запрос";
  if (statusCode === 401) return "Требуется авторизация через Telegram";
  if (statusCode === 403) return "Доступ запрещён";
  if (statusCode === 404) return "Данные не найдены";
  if (statusCode >= 500) return "Ошибка сервера. Попробуйте позже";
  return "Что-то пошло не так. Попробуйте ещё раз";
}

export function messageFromRequestError(err: unknown): string {
  if (err instanceof Error) {
    const match = err.message.match(/^API (\d+): ([\s\S]*)$/);
    if (match) {
      return formatApiError(Number(match[1]), match[2]);
    }
    return err.message;
  }
  return "Что-то пошло не так";
}
