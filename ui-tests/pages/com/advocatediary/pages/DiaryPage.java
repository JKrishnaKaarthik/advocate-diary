package com.advocatediary.pages;

import org.openqa.selenium.By;
import org.openqa.selenium.WebDriver;

public class DiaryPage extends BasePage {
  private final By heading = By.cssSelector("[data-testid='diary-heading']");

  public DiaryPage(WebDriver driver) {
    super(driver);
  }

  public DiaryPage waitUntilLoaded() {
    visible(heading);
    return this;
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
    return By.xpath(
        "//a[contains(@class,'slot') or @data-testid='diary-slot']"
            + "[.//strong[normalize-space()='"
            + clientName
            + "']]"
            + "[.//span[contains(.,'"
            + title
            + "') and contains(.,'"
            + caseNumber
            + "')]]"
            + "[.//span[contains(.,'"
            + court
            + "') and contains(.,'"
            + purpose
            + "')]]");
  }

  public LayoutPage nav() {
    return new LayoutPage(driver);
  }
}
