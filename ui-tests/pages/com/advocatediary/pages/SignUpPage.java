package com.advocatediary.pages;

import org.openqa.selenium.By;
import org.openqa.selenium.TimeoutException;
import org.openqa.selenium.WebDriver;
import org.openqa.selenium.support.ui.ExpectedConditions;

public class SignUpPage extends BasePage {
  private final By heading = testId("auth-heading");
  private final By name = testId("auth-name");
  private final By email = testId("auth-email");
  private final By password = testId("auth-password");
  private final By submit = testId("auth-submit");
  private final By goSignIn = testId("go-sign-in");
  private final By error = By.cssSelector("p.error");

  public SignUpPage(WebDriver driver) {
    super(driver);
  }

  public SignUpPage open(String baseUrl) {
    driver.get(baseUrl + "/sign-up");
    visible(heading);
    return this;
  }

  public DiaryPage createAccount(String lawyerName, String lawyerEmail, String lawyerPassword) {
    type(name, lawyerName);
    type(email, lawyerEmail);
    type(password, lawyerPassword);
    click(submit);
    return new DiaryPage(driver).waitUntilLoaded();
  }

  public SignUpPage submitInvalid(String lawyerName, String lawyerEmail, String lawyerPassword) {
    type(name, lawyerName);
    type(email, lawyerEmail);
    type(password, lawyerPassword);
    click(submit);
    return this;
  }

  public String heading() {
    return textOf(heading);
  }

  public boolean showsPasswordHint() {
    return hintShown("password");
  }

  public String errorText() {
    return textOf(error);
  }

  public SignInPage goToSignIn() {
    click(goSignIn);
    return new SignInPage(driver).waitUntilLoaded();
  }

  public boolean isOpenYourDiary() {
    return currentUrl().contains("/sign-up") && shown(heading);
  }

  public void fillAndCreateAccount(String lawyerName, String lawyerEmail, String lawyerPassword) {
    type(name, lawyerName);
    type(email, lawyerEmail);
    type(password, lawyerPassword);
    pause(1000);
    click(submit);
  }

  public boolean accountCreatedShowsDiary() {
    try {
      wait.until(ExpectedConditions.or(
          ExpectedConditions.visibilityOfElementLocated(testId("diary-heading")),
          ExpectedConditions.visibilityOfElementLocated(error)));
      return shown(testId("diary-heading"));
    } catch (TimeoutException e) {
      return false;
    }
  }

  public boolean hasError() {
    return shown(error);
  }
}
