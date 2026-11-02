package com.advocatediary.tests;

import com.advocatediary.pages.CasePage;
import com.advocatediary.pages.ClientPage;
import com.advocatediary.pages.DiaryPage;
import com.advocatediary.pages.InboxPage;
import com.advocatediary.pages.LayoutPage;
import com.advocatediary.pages.NewClientPage;
import com.advocatediary.pages.SignInPage;
import com.advocatediary.pages.SignUpPage;
import io.github.bonigarcia.wdm.WebDriverManager;
import java.time.Duration;
import java.time.LocalDateTime;
import java.time.ZoneId;
import java.time.format.DateTimeFormatter;
import java.util.HashMap;
import java.util.Map;
import java.util.UUID;
import org.openqa.selenium.WebDriver;
import org.openqa.selenium.chrome.ChromeDriver;
import org.openqa.selenium.chrome.ChromeOptions;
import org.testng.annotations.AfterMethod;
import org.testng.annotations.BeforeMethod;

public abstract class BaseTest {
  protected static final String LAWYER_EMAIL = "bmeherlasya@gmail.com";
  protected static final String LAWYER_PASSWORD = "123@45678";
  protected static final String CLIENT_NAME = "Meera Rao";
  protected static final String CLIENT_PHONE = "9876543210";
  protected static final String CLIENT_EMAIL = "new@example.test";
  protected static final String CLIENT_ADDRESS = "Chennai";
  protected static final String CASE_OPEN_TITLE = "Rao v State";
  protected static final String CASE_TITLE = "Rao v State updated";
  protected static final String CASE_NUMBER = "WP 12 of 2026";
  protected static final String COURT = "Madras High Court";
  protected static final String PURPOSE = "arguments";
  protected static final String OPPOSITE = "State of Tamil Nadu";
  protected static final String CASE_NOTES = "Listed for admission";

  protected WebDriver driver;
  protected String baseUrl;

  protected DiaryPage signInToDiary() {
    return new SignInPage(driver).open(baseUrl).signIn(LAWYER_EMAIL, LAWYER_PASSWORD);
  }

  protected DiaryPage signUpNewAccount() {
    String email = "lawyer+" + UUID.randomUUID() + "@chamber.test";
    return new SignUpPage(driver).open(baseUrl).createAccount("Other Advocate", email, "diary-pass-1");
  }

  protected DiaryPage signupEmptyDiary() {
    DiaryPage diary = signUpNewAccount();
    diary.waitUntilEmpty();
    return diary;
  }

  protected NewClientPage openAddClient() {
    return signInToDiary().nav().openClients().addClient();
  }

  protected ClientPage addClient(String name) {
    return new LayoutPage(driver)
        .openClients()
        .addClient()
        .saveClientWithNotes(name, CLIENT_PHONE, CLIENT_EMAIL, CLIENT_ADDRESS, "First matter");
  }

  protected CasePage openCasePage() {
    signInToDiary();
    return addClient(CLIENT_NAME).openCase(CASE_OPEN_TITLE, COURT, "civil", CASE_NUMBER, OPPOSITE);
  }

  protected CasePage openAndSaveCase(ClientPage client) {
    CasePage matter = client.openCase(CASE_OPEN_TITLE, COURT, "civil", CASE_NUMBER, OPPOSITE);
    matter.fillCaseDetails(
        CASE_TITLE, COURT, CASE_NUMBER, "writ", "active", OPPOSITE, "2026-01-15", CASE_NOTES);
    return matter.clickSaveCaseExpectingTitle(CASE_TITLE);
  }

  protected CasePage saveHearingOnCase(CasePage matter) {
    matter.fillHearing(hearingWhen(), COURT, PURPOSE);
    matter.clickSaveHearing();
    return matter.waitForHearing(COURT, PURPOSE);
  }

  protected DiaryPage completeClientCaseOnDiary() {
    CasePage matter = openAndSaveCase(addClient(CLIENT_NAME));
    saveHearingOnCase(matter);
    DiaryPage diary = matter.nav().openDiary();
    diary.waitForSavedCase(CLIENT_NAME, CASE_TITLE, CASE_NUMBER, COURT, PURPOSE);
    return diary;
  }

  protected boolean diaryShowsMatter(DiaryPage diary) {
    return diary.showsSavedCase(CLIENT_NAME, CASE_TITLE, CASE_NUMBER, COURT, PURPOSE);
  }

  protected boolean inboxShowsMatter(InboxPage inbox) {
    return inbox.showsNextHearing(CLIENT_NAME, CASE_TITLE, CASE_NUMBER);
  }

  protected SignInPage signOut(DiaryPage diary) {
    return diary.nav().signOut();
  }

  private String hearingWhen() {
    return LocalDateTime.now(ZoneId.of("Asia/Kolkata"))
        .plusDays(2)
        .withHour(11)
        .withMinute(15)
        .withSecond(0)
        .withNano(0)
        .format(DateTimeFormatter.ofPattern("yyyy-MM-dd'T'HH:mm"));
  }

  @BeforeMethod
  public void openBrowser() {
    baseUrl = System.getProperty("baseUrl", "http://localhost:8080");
    if (baseUrl == null || baseUrl.isBlank() || baseUrl.contains("${")) {
      baseUrl = "http://localhost:8080";
    }
    WebDriverManager.chromedriver().setup();
    ChromeOptions options = new ChromeOptions();
    boolean headless = Boolean.parseBoolean(System.getProperty("headless", "false"));
    if (headless) {
      options.addArguments("--headless=new", "--disable-dev-shm-usage");
    }
    options.addArguments("--window-size=1400,900", "--disable-gpu", "--no-sandbox", "--start-maximized");
    options.addArguments("--disable-notifications", "--disable-save-password-bubble");
    options.addArguments("--disable-features=PasswordManagerOnboarding,PasswordLeakDetection");
    Map<String, Object> prefs = new HashMap<>();
    prefs.put("credentials_enable_service", false);
    prefs.put("profile.password_manager_enabled", false);
    options.setExperimentalOption("prefs", prefs);
    driver = new ChromeDriver(options);
    driver.manage().timeouts().implicitlyWait(Duration.ofSeconds(0));
    if (!headless) {
      driver.manage().window().maximize();
    }
  }

  @AfterMethod(alwaysRun = true)
  public void closeBrowser() {
    if (driver != null && !Boolean.parseBoolean(System.getProperty("headless", "false"))) {
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
