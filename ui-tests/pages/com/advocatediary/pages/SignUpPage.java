package com.advocatediary.pages;

import org.openqa.selenium.By;
import org.openqa.selenium.TimeoutException;
import org.openqa.selenium.WebDriver;
import org.openqa.selenium.support.ui.ExpectedConditions;

public class SignUpPage extends BasePage {
  private final By heading = By.cssSelector("h1");
  private final By name = By.name("name");
  private final By email = By.name("email");
  private final By password = By.name("password");
  private final By submit = By.cssSelector("[data-testid='auth-submit']");
  private final By goSignIn = By.cssSelector("[data-testid='go-sign-in']");
  private final By error = By.cssSelector("p.error");
  private final By hint = By.cssSelector("p.info");

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

  public String hintText() {
    return textOf(hint);
  }

  public String errorText() {
    return textOf(error);
  }

  public SignInPage goToSignIn() {
    click(goSignIn);
    return new SignInPage(driver).waitUntilLoaded();
  }

  public boolean isOpenYourDiary() {
    return currentUrl().contains("/sign-up") && "Open your diary".equals(heading());
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
          ExpectedConditions.visibilityOfElementLocated(By.cssSelector("[data-testid='diary-heading']")),
          ExpectedConditions.visibilityOfElementLocated(error)));
      return shown(By.cssSelector("[data-testid='diary-heading']"));
    } catch (TimeoutException e) {
      return false;
    }
  }

  public boolean hasError() {
    return shown(error);
  }
}
