package com.advocatediary.tests;

import com.advocatediary.pages.CasePage;
import com.advocatediary.pages.ClientPage;
import com.advocatediary.pages.DiaryPage;
import com.advocatediary.pages.LayoutPage;
import com.advocatediary.pages.SignInPage;
import com.advocatediary.pages.SignUpPage;
import io.github.bonigarcia.wdm.WebDriverManager;
import java.time.Duration;
import java.time.LocalDateTime;
import java.time.ZoneId;
import java.time.format.DateTimeFormatter;
import java.util.UUID;
import org.openqa.selenium.WebDriver;
import org.openqa.selenium.chrome.ChromeDriver;
import org.openqa.selenium.chrome.ChromeOptions;
import org.testng.annotations.AfterMethod;
import org.testng.annotations.BeforeMethod;

public abstract class BaseTest {
  protected static final String LAWYER_NAME = "lasya";
  protected static final String LAWYER_EMAIL = "bmeherlasya@gmail.com";
  protected static final String LAWYER_PASSWORD = "123@45678";
  protected static final String CLIENT_NAME = "Meera Rao";
  protected static final String CASE_TITLE = "Rao v State updated";
  protected static final String CASE_NUMBER = "WP 12 of 2026";
  protected static final String COURT = "Madras High Court";
  protected static final String PURPOSE = "arguments";

  protected WebDriver driver;
  protected String baseUrl;

  protected DiaryPage signInToDiary() {
    return new SignInPage(driver).open(baseUrl).signIn(LAWYER_EMAIL, LAWYER_PASSWORD);
  }

  protected DiaryPage signUpNewAccount() {
    String email = "lawyer+" + UUID.randomUUID() + "@chamber.test";
    return new SignUpPage(driver).open(baseUrl).createAccount("Other Advocate", email, "diary-pass-1");
  }

  protected ClientPage addClient(String name) {
    return new LayoutPage(driver)
        .openClients()
        .addClient()
        .saveClientWithNotes(name, "9876543210", "new@example.test", "Chennai", "First matter");
  }

  protected CasePage openAndSaveCase(ClientPage client) {
    CasePage matter = client.openCase(
        "Rao v State",
        COURT,
        "civil",
        CASE_NUMBER,
        "State of Tamil Nadu");
    matter.fillCaseDetails(
        CASE_TITLE,
        COURT,
        CASE_NUMBER,
        "writ",
        "active",
        "State of Tamil Nadu",
        "2026-01-15",
        "Listed for admission");
    return matter.clickSaveCaseExpectingTitle(CASE_TITLE);
  }

  protected CasePage saveHearingOnCase(CasePage matter) {
    String when = LocalDateTime.now(ZoneId.of("Asia/Kolkata"))
        .plusDays(2)
        .withHour(11)
        .withMinute(15)
        .withSecond(0)
        .withNano(0)
        .format(DateTimeFormatter.ofPattern("yyyy-MM-dd'T'HH:mm"));
    matter.fillHearing(when, COURT, PURPOSE);
    matter.clickSaveHearing();
    return matter.waitForHearing(COURT, PURPOSE);
  }

  @BeforeMethod
  public void openBrowser() {
    baseUrl = System.getProperty("baseUrl", "http://localhost:8080");
    if (baseUrl == null || baseUrl.isBlank() || baseUrl.contains("${")) {
      baseUrl = "http://localhost:8080";
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
        Thread.sleep(6000);
      } catch (InterruptedException e) {
        Thread.currentThread().interrupt();
      }
    }
    if (driver != null) {
      driver.quit();
    }
  }
}
