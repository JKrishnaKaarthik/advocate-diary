package com.advocatediary.pages;

import org.openqa.selenium.By;
import org.openqa.selenium.WebDriver;

public class DiaryPage extends BasePage {
  private final By heading = testId("diary-heading");

  public DiaryPage(WebDriver driver) {
    super(driver);
  }

  public DiaryPage waitUntilLoaded() {
    visible(heading);
    return this;
  }

  public boolean isOpen() {
    return shown(heading);
  }

  public DiaryPage waitUntilEmpty() {
    visible(empty());
    return this;
  }

  public boolean isEmpty() {
    return shown(empty()) && !shown(testId("diary-slot"));
  }

  private By empty() {
    return By.cssSelector("[data-testid='diary-empty'], div.empty");
  }

  public String heading() {
    return textOf(heading);
  }

  public boolean showsSavedCase(String clientName, String title, String caseNumber, String court, String purpose) {
    return shown(slot(clientName, title, caseNumber, court, purpose));
  }

  public DiaryPage waitForSavedCase(String clientName, String title, String caseNumber, String court, String purpose) {
    visible(slot(clientName, title, caseNumber, court, purpose));
    return this;
  }

  private By slot(String clientName, String title, String caseNumber, String court, String purpose) {
    return marked(
        "diary-slot",
        "data-client", clientName,
        "data-title", title,
        "data-case-number", caseNumber,
        "data-court", court,
        "data-purpose", purpose);
  }

  public LayoutPage nav() {
    return new LayoutPage(driver);
  }
}
