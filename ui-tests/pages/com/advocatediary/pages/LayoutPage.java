package com.advocatediary.pages;

import org.openqa.selenium.By;
import org.openqa.selenium.WebDriver;

public class LayoutPage extends BasePage {
  private final By diary = By.cssSelector("[data-testid='nav-diary']");
  private final By clients = By.cssSelector("[data-testid='nav-clients']");
  private final By inbox = By.cssSelector("[data-testid='nav-inbox']");
  private final By signOut = By.cssSelector("[data-testid='sign-out']");

  public LayoutPage(WebDriver driver) {
    super(driver);
  }

  public DiaryPage openDiary() {
    click(diary);
    return new DiaryPage(driver).waitUntilLoaded();
  }

  public ClientsPage openClients() {
    click(clients);
    return new ClientsPage(driver).waitUntilLoaded();
  }

  public InboxPage openInbox() {
    click(inbox);
    return new InboxPage(driver).waitUntilLoaded();
  }

  public SignInPage signOut() {
    click(signOut);
    return new SignInPage(driver).waitUntilLoaded();
  }
}
