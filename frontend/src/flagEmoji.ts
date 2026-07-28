// Converts a 2-letter ISO country code into its flag emoji using the
// regional indicator symbol Unicode trick (each letter A-Z maps to the
// codepoint range U+1F1E6-U+1F1FF).
export function isoToFlagEmoji(iso: string): string {
  if (!iso || iso.length !== 2) return "";
  const codePoints = [...iso.toUpperCase()].map((c) => 127397 + c.charCodeAt(0));
  return String.fromCodePoint(...codePoints);
}
