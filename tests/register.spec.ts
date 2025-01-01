import { test, expect } from "@playwright/test";
import { generateEmail, generateRandomName } from "./utils/user";

test("sign up button works", async ({ page }) => {
    await page.goto("/");
    await page.getByRole("link", { name: "Sign up" }).click();
    await expect(page).toHaveTitle(/Register/);
    await expect(page).toHaveURL(/register/);
});


test("register form works", async ({ page }) => {
    await page.goto("/");
    await page.getByRole("link", { name: "Sign up" }).click();

    // fill out register form
    const nickname = generateRandomName();
    const email = generateEmail(nickname);
    const password = "VerySecurePassword1234!@";
    await page.getByPlaceholder("Nickname").fill(nickname);
    await page.getByPlaceholder("Email").fill(email);
    await page.locator("//input[@id='password']").fill(password);
    await page.locator("//input[@id='confirm_password']").fill(password);

    // click on 'register' button
    await page.getByRole("button", { name: "Register" }).click();

    // expect success
    await expect(page).toHaveURL(/login/);
    await expect(page.getByText(/Account created successfully !/)).toBeVisible();
});
