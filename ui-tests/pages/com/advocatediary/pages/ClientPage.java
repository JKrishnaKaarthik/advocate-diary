package com.advocatediary.pages;

import org.openqa.selenium.By;
import org.openqa.selenium.WebDriver;
import org.openqa.selenium.support.ui.Select;

public class ClientPage extends BasePage {
  private final By heading = testId("client-heading");
  private final By openCaseHeading = testId("open-case-heading");
  private final By title = testId("open-case-title");
  private final By court = testId("open-case-court");
  private final By caseType = testId("open-case-type");
  private final By caseNumber = testId("open-case-number");
  private final By oppositeParty = testId("open-case-opposite");
  private final By openCase = testId("open-case");

  public ClientPage(WebDriver driver) {
    super(driver);
  }

  public ClientPage waitUntilLoaded() {
    visible(heading);
    return this;
  }

  public String heading() {
    return textOf(heading);
  }

  public boolean showsOpenACase() {
    return shown(openCaseHeading) && shown(title) && shown(court) && shown(openCase);
  }

  public void fillOpenCase(String caseTitle, String courtName, String type, String number, String opposite) {
    type(title, caseTitle);
    type(court, courtName);
    new Select(visible(caseType)).selectByVisibleText(type);
    type(caseNumber, number);
    type(oppositeParty, opposite);
    pause(600);
  }

  public CasePage clickOpenCase() {
    click(openCase);
    return new CasePage(driver).waitUntilLoaded();
  }

  public CasePage openCase(String caseTitle, String courtName, String type, String number, String opposite) {
    fillOpenCase(caseTitle, courtName, type, number, opposite);
    return clickOpenCase();
  }
}
