"use client";

import { useCallback, useEffect, useState } from "react";
import { fetchContent } from "../lib/api/display-hello-word-from-database";
import styles from "./DisplayHelloWordFromDatabase.module.css";

type Status = "loading" | "loaded" | "error";

export default function DisplayHelloWordFromDatabase() {
  const [status, setStatus] = useState<Status>("loading");
  const [value, setValue] = useState("");

  const load = useCallback(async () => {
    setStatus("loading");
    try {
      const data = await fetchContent();
      if (!data.value) {
        setStatus("error");
        return;
      }
      setValue(data.value);
      setStatus("loaded");
    } catch {
      setStatus("error");
    }
  }, []);

  useEffect(() => {
    load();
  }, [load]);

  const is = (s: Status) => (status === s ? "" : styles.hidden);

  return (
    <div className={styles.panel} role="status" aria-live="polite">
      <p className={`${styles.loading} ${is("loading")}`}>
        Loading
        <span className={styles.dot}>.</span>
        <span className={styles.dot}>.</span>
        <span className={styles.dot}>.</span>
      </p>
      <h1 className={`${styles.heading} ${is("loaded")}`}>{value}</h1>
      <p className={`${styles.error} ${is("error")}`} role="alert">
        Could not load the text from the database.
      </p>
      {status === "error" && (
        <button className={styles.button} type="button" onClick={load}>
          Retry
        </button>
      )}
    </div>
  );
}
