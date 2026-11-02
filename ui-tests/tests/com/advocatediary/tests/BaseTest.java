package com.advocatediary.tests;

import com.advocatediary.pages.DiaryPage;
import com.advocatediary.pages.SignInPage;
import io.github.bonigarcia.wdm.WebDriverManager;
import java.time.Duration;
import org.openqa.selenium.WebDriver;
import org.openqa.selenium.chrome.ChromeDriver;
import org.openqa.selenium.chrome.ChromeOptions;
import org.testng.annotations.AfterMethod;
import org.testng.annotations.BeforeMethod;

public abstract class BaseTest {
  protected static final String LAWYER_NAME = "lasya";
  protected static final String LAWYER_EMAIL = "bmeherlasya@gmail.com";
  protected static final String LAWYER_PASSWORD = "123@45678";

  protected WebDriver driver;
  protected String baseUrl;

  protected DiaryPage signInToDiary() {
    return new SignInPage(driver).open(baseUrl).signIn(LAWYER_EMAIL, LAWYER_PASSWORD);
  }

  @BeforeMethod
  public void openBrowser() {
    baseUrl = System.getProperty("baseUrl", "http://localhost:5173");
    if (baseUrl == null || baseUrl.isBlank() || baseUrl.contains("${")) {
      baseUrl = "http://localhost:5173";
    }
    WebDriverManager.chromedriver().setup();
    ChromeOptions options = new ChromeOptions();
    boolean headless = Boolean.parseBoolean(System.getProperty("headless", "true"));
    if (headless) {
      options.addArguments("--headless=new");
    }
    options.addArguments("--window-size=1400,900", "--disable-gpu", "--no-sandbox");
    options.addArguments("--start-maximized");
    driver = new ChromeDriver(options);
    driver.manage().timeouts().implicitlyWait(Duration.ofSeconds(0));
    if (!headless) {
      driver.manage().window().maximize();
    }
  }

  @AfterMethod(alwaysRun = true)
  public void closeBrowser() {
    if (driver != null && !Boolean.parseBoolean(System.getProperty("headless", "true"))) {
      try {
        Thread.sleep(2500);
      } catch (InterruptedException e) {
        Thread.currentThread().interrupt();
      }
    }
    if (driver != null) {
      driver.quit();
    }
  }
}
