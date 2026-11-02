package com.advocatediary.pages;

import java.nio.file.Path;
import org.openqa.selenium.By;
import org.openqa.selenium.JavascriptExecutor;
import org.openqa.selenium.WebDriver;
import org.openqa.selenium.support.ui.ExpectedConditions;
import org.openqa.selenium.support.ui.Select;

public class CasePage extends BasePage {
  private final By heading = testId("case-heading");
  private final By details = testId("case-details");
  private final By title = testId("case-title");
  private final By court = testId("case-court");
  private final By caseNumber = testId("case-number");
  private final By caseType = testId("case-type");
  private final By status = testId("case-status");
  private final By oppositeParty = testId("case-opposite");
  private final By filingDate = testId("case-filing-date");
  private final By notes = testId("case-notes");
  private final By saveCaseButton = testId("save-case");
  private final By fileInput = testId("upload-file");
  private final By labelSelect = testId("upload-label");
  private final By uploadButton = testId("upload");
  private final By fixHearingHeading = testId("fix-hearing-heading");
  private final By hearingWhen = testId("hearing-when");
  private final By hearingCourt = testId("hearing-court");
  private final By hearingPurpose = testId("hearing-purpose");
  private final By saveHearingButton = testId("save-hearing");

  public CasePage(WebDriver driver) {
    super(driver);
  }

  public CasePage waitUntilLoaded() {
    visible(heading);
    visible(details);
    return this;
  }

  public LayoutPage nav() {
    return new LayoutPage(driver);
  }

  public String heading() {
    return textOf(heading);
  }

  public boolean showsCaseDetails() {
    return shown(details) && shown(title) && shown(court) && shown(saveCaseButton);
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
    return shown(fileInput) && shown(labelSelect) && shown(uploadButton);
  }

  public boolean acceptsPdfOnly() {
    String accept = visible(fileInput).getAttribute("accept");
    return accept != null && accept.contains("pdf") && !accept.contains(".txt") && !accept.contains(".doc");
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

  public boolean showsFileRequired() {
    return hintShown("file");
  }

  public boolean showsPaper(String fileName, String label) {
    return shown(paper(fileName, label));
  }

  public CasePage waitForPaper(String fileName, String label) {
    visible(paper(fileName, label));
    return this;
  }

  private By paper(String fileName, String label) {
    return marked("paper", "data-name", fileName, "data-label", label);
  }

  public boolean showsFixAHearing() {
    return shown(fixHearingHeading)
        && shown(hearingWhen)
        && shown(hearingCourt)
        && shown(hearingPurpose)
        && shown(saveHearingButton);
  }

  public CasePage fillHearing(String when, String courtName, String purpose) {
    ((JavascriptExecutor) driver).executeScript(
        "arguments[0].value = arguments[1];"
            + " arguments[0].dispatchEvent(new Event('input', { bubbles: true }));"
            + " arguments[0].dispatchEvent(new Event('change', { bubbles: true }));",
        visible(hearingWhen),
        when);
    type(hearingCourt, courtName);
    new Select(visible(hearingPurpose)).selectByVisibleText(purpose);
    pause(600);
    return this;
  }

  public CasePage clickSaveHearing() {
    click(saveHearingButton);
    pause(600);
    return this;
  }

  public boolean showsDateRequired() {
    return hintShown("scheduled_at");
  }

  public boolean showsHearing(String courtName, String purpose) {
    return shown(hearingRow(courtName, purpose));
  }

  public CasePage waitForHearing(String courtName, String purpose) {
    visible(hearingRow(courtName, purpose));
    return this;
  }

  private By hearingRow(String courtName, String purpose) {
    return marked("hearing", "data-court", courtName, "data-purpose", purpose, "data-status", "scheduled");
  }
}
