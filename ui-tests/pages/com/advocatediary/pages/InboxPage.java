package com.advocatediary.pages;

import org.openqa.selenium.By;
import org.openqa.selenium.WebDriver;

public class InboxPage extends BasePage {
  private final By heading = By.cssSelector("[data-testid='inbox-heading'], h1");

  public InboxPage(WebDriver driver) {
    super(driver);
  }

  public InboxPage waitUntilLoaded() {
    visible(heading);
    return this;
  }

  public String heading() {
    return textOf(heading);
  }

  public boolean showsNextHearing(String clientName, String title, String caseNumber) {
    return shown(card(clientName, title, caseNumber));
  }

  public InboxPage waitForNextHearing(String clientName, String title, String caseNumber) {
    visible(card(clientName, title, caseNumber));
    return this;
  }

  private By card(String clientName, String title, String caseNumber) {
    return By.xpath(
        "//article[contains(@class,'hearing-card') or @data-testid='inbox-hearing']"
            + "[.//h3[normalize-space()='"
            + clientName
            + "']]"
            + "[.//p[contains(@class,'matter') and contains(.,'"
            + title
            + "') and contains(.,'"
            + caseNumber
            + "')]]");
  }
}
