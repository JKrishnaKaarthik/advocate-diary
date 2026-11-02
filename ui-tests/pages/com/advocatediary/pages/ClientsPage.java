package com.advocatediary.pages;

import org.openqa.selenium.By;
import org.openqa.selenium.WebDriver;

public class ClientsPage extends BasePage {
  private final By heading = By.cssSelector("[data-testid='clients-heading']");
  private final By addClient = By.cssSelector("[data-testid='add-client']");

  public ClientsPage(WebDriver driver) {
    super(driver);
  }

  public ClientsPage waitUntilLoaded() {
    visible(heading);
    return this;
  }

  public String heading() {
    return textOf(heading);
  }

  public NewClientPage addClient() {
    click(addClient);
    return new NewClientPage(driver).waitUntilLoaded();
  }

  public boolean isClientsPage() {
    return currentUrl().contains("/clients")
        && !currentUrl().contains("/clients/new")
        && "Clients".equals(heading());
  }

  public boolean hasClientNamed(String name) {
    return driver.getPageSource().contains(name);
  }
}
