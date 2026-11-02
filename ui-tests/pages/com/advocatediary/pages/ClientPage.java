package com.advocatediary.pages;

import org.openqa.selenium.By;
import org.openqa.selenium.WebDriver;
import org.openqa.selenium.support.ui.Select;

public class ClientPage extends BasePage {
  private final By heading = By.cssSelector("[data-testid='client-heading']");
  private final By openCaseHeading = By.cssSelector("[data-testid='open-case-heading'], h2");
  private final By title = By.name("title");
  private final By court = By.name("court_name");
  private final By caseType = By.name("case_type");
  private final By caseNumber = By.name("case_number");
  private final By oppositeParty = By.name("opposite_party");
  private final By openCase = By.cssSelector("[data-testid='open-case']");

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
    return driver.getPageSource().contains("Open a case")
        && shown(title)
        && shown(court)
        && shown(openCase);
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
