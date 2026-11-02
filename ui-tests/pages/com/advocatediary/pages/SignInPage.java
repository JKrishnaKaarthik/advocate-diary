package com.advocatediary.pages;

import org.openqa.selenium.By;
import org.openqa.selenium.TimeoutException;
import org.openqa.selenium.WebDriver;
import org.openqa.selenium.support.ui.ExpectedConditions;

public class SignInPage extends BasePage {
  private final By heading = testId("auth-heading");
  private final By email = testId("auth-email");
  private final By password = testId("auth-password");
  private final By submit = testId("auth-submit");
  private final By goSignUp = testId("go-sign-up");

  public SignInPage(WebDriver driver) {
    super(driver);
  }

  public SignInPage open(String baseUrl) {
    driver.get(baseUrl + "/sign-in");
    visible(heading);
    pause(800);
    return this;
  }

  public SignInPage waitUntilLoaded() {
    visible(heading);
    return this;
  }

  public DiaryPage signIn(String lawyerEmail, String lawyerPassword) {
    type(email, lawyerEmail);
    type(password, lawyerPassword);
    pause(800);
    click(submit);
    return new DiaryPage(driver).waitUntilLoaded();
  }

  public String heading() {
    return textOf(heading);
  }

  public SignUpPage clickCreateAnAccount() {
    pause(600);
    click(goSignUp);
    try {
      wait.until(ExpectedConditions.urlContains("/sign-up"));
    } catch (TimeoutException ignored) {
      // stay on /sign-in when the create-account page does not open
    }
    return new SignUpPage(driver);
  }

  public boolean isOnSignIn() {
    return currentUrl().contains("/sign-in") && shown(heading);
  }
}
