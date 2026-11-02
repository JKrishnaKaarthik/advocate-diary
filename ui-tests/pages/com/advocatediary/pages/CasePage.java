package com.advocatediary.pages;

import java.nio.file.Path;
import org.openqa.selenium.By;
import org.openqa.selenium.JavascriptExecutor;
import org.openqa.selenium.WebDriver;
import org.openqa.selenium.support.ui.ExpectedConditions;
import org.openqa.selenium.support.ui.Select;

public class CasePage extends BasePage {
  private final By heading = By.cssSelector("[data-testid='case-heading']");
  private final By title = By.name("title");
  private final By court = By.name("court_name");
  private final By caseNumber = By.name("case_number");
  private final By caseType = By.name("case_type");
  private final By status = By.name("status");
  private final By oppositeParty = By.name("opposite_party");
  private final By filingDate = By.name("filing_date");
  private final By notes = By.name("notes");
  private final By detailsHeading = By.xpath("//h2[normalize-space()='Case details']");
  private final By saveCaseButton = By.xpath("//button[normalize-space()='Save case']");
  private final By papersHeading = By.xpath("//h2[normalize-space()='Papers']");
  private final By fileInput = By.cssSelector("[data-testid='upload-file'], input[name='file'][type='file']");
  private final By labelSelect = By.cssSelector("[data-testid='upload-label'], select[name='label']");
  private final By uploadButton = By.xpath("//button[@data-testid='upload' or normalize-space()='Upload']");
  private final By uploadHint = By.xpath("//form[.//button[@data-testid='upload' or normalize-space()='Upload']]//p[contains(@class,'info')]");
  private final By fixHearingHeading = By.xpath("//h2[@data-testid='fix-hearing-heading' or normalize-space()='Fix a hearing']");
  private final By hearingWhen = By.cssSelector("[data-testid='hearing-when'], input[name='scheduled_at'][type='datetime-local']");
  private final By hearingCourt = By.cssSelector("[data-testid='hearing-court']");
  private final By hearingCourtFallback = By.xpath("//h2[normalize-space()='Fix a hearing']/following::input[@name='court_name'][1]");
  private final By hearingPurpose = By.cssSelector("[data-testid='hearing-purpose']");
  private final By hearingPurposeFallback = By.xpath("//h2[normalize-space()='Fix a hearing']/following::select[@name='purpose'][1]");
  private final By saveHearingButton = By.xpath("//button[@data-testid='save-hearing' or normalize-space()='Save hearing']");
  private final By hearingSavedNote = By.cssSelector("p.note");
  private final By hearingWhenHint = By.xpath("//form[.//button[@data-testid='save-hearing' or normalize-space()='Save hearing']]//p[contains(@class,'info')]");

  public CasePage(WebDriver driver) {
    super(driver);
  }

  public CasePage waitUntilLoaded() {
    visible(heading);
    visible(detailsHeading);
    return this;
  }

  public LayoutPage nav() {
    return new LayoutPage(driver);
  }

  public String heading() {
    return textOf(heading);
  }

  public boolean showsCaseDetails() {
    return shown(detailsHeading) && shown(title) && shown(court) && shown(saveCaseButton);
  }

  public void fillCaseDetails(
      String caseTitle,
      String courtName,
      String number,
      String type,
      String caseStatus,
      String opposite,
      String filedOn,
      String note) {
    type(title, caseTitle);
    type(court, courtName);
    type(caseNumber, number);
    new Select(visible(caseType)).selectByVisibleText(type);
    new Select(visible(status)).selectByVisibleText(caseStatus);
    type(oppositeParty, opposite);
    ((JavascriptExecutor) driver).executeScript(
        "arguments[0].value = arguments[1]; arguments[0].dispatchEvent(new Event('input', { bubbles: true }));",
        visible(filingDate),
        filedOn);
    type(notes, note);
    pause(800);
  }

  public CasePage clickSaveCase() {
    click(saveCaseButton);
    pause(500);
    waitUntilLoaded();
    return this;
  }

  public CasePage clickSaveCaseExpectingTitle(String expectedTitle) {
    click(saveCaseButton);
    wait.until(ExpectedConditions.textToBe(heading, expectedTitle));
    return this;
  }

  public String titleValue() {
    return visible(title).getAttribute("value");
  }

  public String notesValue() {
    return visible(notes).getAttribute("value");
  }

  public boolean showsFileUpload() {
    return shown(papersHeading) && shown(fileInput) && shown(labelSelect) && shown(uploadButton);
  }

  public CasePage chooseFile(Path file) {
    visible(fileInput).sendKeys(file.toAbsolutePath().toString());
    pause(600);
    return this;
  }

  public CasePage selectLabel(String label) {
    new Select(visible(labelSelect)).selectByVisibleText(label);
    pause(400);
    return this;
  }

  public CasePage clickUpload() {
    click(uploadButton);
    pause(600);
    return this;
  }

  public boolean showsChooseAFileHint() {
    return shown(uploadHint) && textOf(uploadHint).contains("Choose a file to upload");
  }

  public boolean showsPaper(String fileName, String label) {
    return shown(paper(fileName, label));
  }

  public CasePage waitForPaper(String fileName, String label) {
    visible(paper(fileName, label));
    return this;
  }

  private By paper(String fileName, String label) {
    return By.xpath(
        "//li[(contains(@class,'doc') or @data-testid='paper')]"
            + "[.//strong[normalize-space()='"
            + fileName
            + "']]"
            + "[.//span[contains(.,'"
            + label
            + "')]]");
  }

  public boolean showsFixAHearing() {
    return shown(fixHearingHeading)
        && shown(hearingWhen)
        && shown(hearingCourtField())
        && shown(hearingPurposeField())
        && shown(saveHearingButton);
  }

  public CasePage fillHearing(String when, String courtName, String purpose) {
    ((JavascriptExecutor) driver).executeScript(
        "arguments[0].value = arguments[1];"
            + " arguments[0].dispatchEvent(new Event('input', { bubbles: true }));"
            + " arguments[0].dispatchEvent(new Event('change', { bubbles: true }));",
        visible(hearingWhen),
        when);
    type(hearingCourtField(), courtName);
    new Select(visible(hearingPurposeField())).selectByVisibleText(purpose);
    pause(600);
    return this;
  }

  public CasePage clickSaveHearing() {
    click(saveHearingButton);
    pause(600);
    return this;
  }

  public CasePage clickSaveHearingExpectingSaved() {
    click(saveHearingButton);
    wait.until(ExpectedConditions.textToBePresentInElementLocated(hearingSavedNote, "Hearing saved."));
    return this;
  }

  public boolean showsDateAndTimeRequired() {
    return shown(hearingWhenHint) && textOf(hearingWhenHint).contains("Date and time is required");
  }

  public boolean showsHearing(String courtName, String purpose) {
    return shown(hearingRow(courtName, purpose));
  }

  public CasePage waitForHearing(String courtName, String purpose) {
    visible(hearingRow(courtName, purpose));
    return this;
  }

  public String hearingSavedNote() {
    return textOf(hearingSavedNote);
  }

  private By hearingCourtField() {
    return shown(hearingCourt) ? hearingCourt : hearingCourtFallback;
  }

  private By hearingPurposeField() {
    return shown(hearingPurpose) ? hearingPurpose : hearingPurposeFallback;
  }

  private By hearingRow(String courtName, String purpose) {
    return By.xpath(
        "//li[@data-testid='hearing' or .//span[contains(.,' · ')]]"
            + "[.//span[contains(.,'"
            + courtName
            + "') and contains(.,'"
            + purpose
            + "')]]"
            + "[.//strong[contains(.,'scheduled')]]");
  }
}
