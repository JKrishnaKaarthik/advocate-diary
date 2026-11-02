package com.advocatediary.pages;

import org.openqa.selenium.By;
import org.openqa.selenium.WebDriver;

public class InboxPage extends BasePage {
  private final By heading = testId("inbox-heading");

  public InboxPage(WebDriver driver) {
    super(driver);
  }

  public InboxPage waitUntilLoaded() {
    visible(heading);
    return this;
  }

  public boolean isOpen() {
    return shown(heading);
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
    return marked(
        "inbox-hearing",
        "data-client", clientName,
        "data-title", title,
        "data-case-number", caseNumber);
  }
}
