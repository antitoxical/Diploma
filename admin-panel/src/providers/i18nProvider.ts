import type { I18nProvider } from '@refinedev/core';
import { ruTranslations } from './locales/ru';

export const i18nProvider: I18nProvider = {
  translate: (key, options?, defaultMessage?) => {
    const mapped = ruTranslations[key];
    if (mapped) return mapped;

    // Refine часто вызывает translate(key, defaultMessage) — второй аргумент это текст по умолчанию
    if (typeof options === 'string' && defaultMessage === undefined) {
      return options;
    }

    return defaultMessage ?? key;
  },
  changeLocale: async () => {},
  getLocale: () => 'ru',
};
