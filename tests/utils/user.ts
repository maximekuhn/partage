export function generateRandomName(): string {
    const baseName = "johndoe";
    const randomNumber = Math.floor(10_000 + Math.random() * 90_000).toString();
    return baseName + randomNumber;
}

export function generateEmail(name: string): string {
    return name + "@gmail.com";
}
