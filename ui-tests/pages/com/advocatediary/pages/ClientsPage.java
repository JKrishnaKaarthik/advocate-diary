package com.advocatediary.pages;

import org.openqa.selenium.By;
import org.openqa.selenium.WebDriver;

public class ClientsPage extends BasePage {
  private final By heading = testId("clients-heading");
  private final By addClient = testId("add-client");

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
        && shown(heading);
  }

  public ClientsPage waitUntilEmpty() {
    visible(By.cssSelector("[data-testid='clients-empty'], p.muted"));
    return this;
  }

  public boolean isEmpty() {
    return shown(By.cssSelector("[data-testid='clients-empty'], p.muted")) && !shown(testId("client-row"));
  }

  public boolean hasClient(String name) {
    return shown(marked("client-row", "data-name", name))
        || driver.findElements(By.xpath("//strong[normalize-space()='" + name + "']"))
            .stream()
            .anyMatch(el -> el.isDisplayed());
  }
}
