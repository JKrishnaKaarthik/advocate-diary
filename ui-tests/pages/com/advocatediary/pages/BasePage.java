package com.advocatediary.pages;

import java.time.Duration;
import org.openqa.selenium.By;
import org.openqa.selenium.Keys;
import org.openqa.selenium.WebDriver;
import org.openqa.selenium.WebElement;
import org.openqa.selenium.support.ui.ExpectedConditions;
import org.openqa.selenium.support.ui.WebDriverWait;

public abstract class BasePage {
  protected final WebDriver driver;
  protected final WebDriverWait wait;

  protected BasePage(WebDriver driver) {
    this.driver = driver;
    this.wait = new WebDriverWait(driver, Duration.ofSeconds(15));
  }

  protected WebElement visible(By locator) {
    return wait.until(ExpectedConditions.visibilityOfElementLocated(locator));
  }

  protected WebElement clickable(By locator) {
    return wait.until(ExpectedConditions.elementToBeClickable(locator));
  }

  protected void type(By locator, String value) {
    WebElement field = visible(locator);
    field.click();
    field.sendKeys(Keys.chord(Keys.CONTROL, "a"));
    field.sendKeys(Keys.DELETE);
    if (watchMode()) {
      for (char c : value.toCharArray()) {
        field.sendKeys(String.valueOf(c));
        pause(90);
      }
    } else {
      field.sendKeys(value);
    }
  }

  protected boolean watchMode() {
    return !Boolean.parseBoolean(System.getProperty("headless", "true"));
  }

  protected void pause(long millis) {
    try {
      Thread.sleep(millis);
    } catch (InterruptedException e) {
      Thread.currentThread().interrupt();
    }
  }

  protected void click(By locator) {
    clickable(locator).click();
  }

  protected String textOf(By locator) {
    return visible(locator).getText();
  }

  protected boolean shown(By locator) {
    return driver.findElements(locator).stream().anyMatch(WebElement::isDisplayed);
  }

  public String currentUrl() {
    return driver.getCurrentUrl();
  }

  protected static By testId(String id) {
    return By.cssSelector("[data-testid='" + id + "']");
  }

  protected static By marked(String id, String... nameValues) {
    StringBuilder css = new StringBuilder("[data-testid='").append(id).append("']");
    for (int i = 0; i + 1 < nameValues.length; i += 2) {
      css.append('[')
          .append(nameValues[i])
          .append("='")
          .append(nameValues[i + 1].replace("\\", "\\\\").replace("'", "\\'"))
          .append("']");
    }
    return By.cssSelector(css.toString());
  }

  protected boolean hintShown(String field) {
    return shown(testId("hint-" + field));
  }
}
