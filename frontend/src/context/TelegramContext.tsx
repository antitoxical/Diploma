import React, { createContext, useContext, useEffect, useState, ReactNode } from "react";
import { api } from "../api/client";

export interface TgUser {
  id: number;
  first_name: string;
  last_name?: string;
  username?: string;
  language_code?: string;
  photo_url?: string;
}

interface TelegramContextType {
  user: TgUser | null;
  initData: string;
  ready: boolean;
}

const TelegramContext = createContext<TelegramContextType>({
  user: null,
  initData: "",
  ready: false,
});

export const TelegramProvider = ({ children }: { children: ReactNode }) => {
  const [user, setUser] = useState<TgUser | null>(null);
  const [initData, setInitData] = useState("");
  const [ready, setReady] = useState(false);

  useEffect(() => {
    let tgUser: TgUser | null = null;
    try {
      const tg = (window as any).Telegram?.WebApp;
      if (tg) {
        tg.ready();
        tg.expand();
        setInitData(tg.initData || "");
        if (tg.initDataUnsafe?.user) {
          tgUser = tg.initDataUnsafe.user;
          setUser(tgUser);
        }
      }
    } catch {
      // not in Telegram context
    }
    setReady(true);

    if (tgUser) {
      api.getUser(tgUser.id).catch(() => {
        if (!tgUser) return;
        api.createUser({
          telegram_id: tgUser.id,
          username: tgUser.username,
          first_name: tgUser.first_name,
          last_name: tgUser.last_name,
        }).catch(() => {});
      });
    }
  }, []);

  return (
    <TelegramContext.Provider value={{ user, initData, ready }}>
      {children}
    </TelegramContext.Provider>
  );
};

export const useTelegram = () => useContext(TelegramContext);
